package hosts

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
	"github.com/coredns/coredns/plugin"
)

func parseLiteralIP(addr string) net.IP {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if i := strings.Index(addr, "%"); i >= 0 {
		addr = addr[0:i]
	}
	return net.ParseIP(addr)
}
func absDomainName(b string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return plugin.Name(b).Normalize()
}

type hostsMap struct {
	byNameV4	map[string][]net.IP
	byNameV6	map[string][]net.IP
	byAddr		map[string][]string
}

func newHostsMap() *hostsMap {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &hostsMap{byNameV4: make(map[string][]net.IP), byNameV6: make(map[string][]net.IP), byAddr: make(map[string][]string)}
}
func (h *hostsMap) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := 0
	for _, v4 := range h.byNameV4 {
		l += len(v4)
	}
	for _, v6 := range h.byNameV6 {
		l += len(v6)
	}
	for _, a := range h.byAddr {
		l += len(a)
	}
	return l
}

type Hostsfile struct {
	sync.RWMutex
	Origins	[]string
	hmap	*hostsMap
	inline	*hostsMap
	path	string
	mtime	time.Time
	size	int64
}

func (h *Hostsfile) readHosts() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	file, err := os.Open(h.path)
	if err != nil {
		return
	}
	defer file.Close()
	stat, err := file.Stat()
	if err == nil && h.mtime.Equal(stat.ModTime()) && h.size == stat.Size() {
		return
	}
	newMap := h.parse(file, h.inline)
	log.Debugf("Parsed hosts file into %d entries", newMap.Len())
	h.Lock()
	h.hmap = newMap
	h.mtime = stat.ModTime()
	h.size = stat.Size()
	h.Unlock()
}
func (h *Hostsfile) initInline(inline []string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(inline) == 0 {
		return
	}
	hmap := newHostsMap()
	h.inline = h.parse(strings.NewReader(strings.Join(inline, "\n")), hmap)
	*h.hmap = *h.inline
}
func (h *Hostsfile) parse(r io.Reader, override *hostsMap) *hostsMap {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	hmap := newHostsMap()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		if i := bytes.Index(line, []byte{'#'}); i >= 0 {
			line = line[0:i]
		}
		f := bytes.Fields(line)
		if len(f) < 2 {
			continue
		}
		addr := parseLiteralIP(string(f[0]))
		if addr == nil {
			continue
		}
		ver := ipVersion(string(f[0]))
		for i := 1; i < len(f); i++ {
			name := absDomainName(string(f[i]))
			if plugin.Zones(h.Origins).Matches(name) == "" {
				continue
			}
			switch ver {
			case 4:
				hmap.byNameV4[name] = append(hmap.byNameV4[name], addr)
			case 6:
				hmap.byNameV6[name] = append(hmap.byNameV6[name], addr)
			default:
				continue
			}
			hmap.byAddr[addr.String()] = append(hmap.byAddr[addr.String()], name)
		}
	}
	if override == nil {
		return hmap
	}
	for name := range override.byNameV4 {
		hmap.byNameV4[name] = append(hmap.byNameV4[name], override.byNameV4[name]...)
	}
	for name := range override.byNameV4 {
		hmap.byNameV6[name] = append(hmap.byNameV6[name], override.byNameV6[name]...)
	}
	for addr := range override.byAddr {
		hmap.byAddr[addr] = append(hmap.byAddr[addr], override.byAddr[addr]...)
	}
	return hmap
}
func ipVersion(s string) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return 4
		case ':':
			return 6
		}
	}
	return 0
}
func (h *Hostsfile) LookupStaticHostV4(host string) []net.IP {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h.RLock()
	defer h.RUnlock()
	if len(h.hmap.byNameV4) != 0 {
		if ips, ok := h.hmap.byNameV4[absDomainName(host)]; ok {
			ipsCp := make([]net.IP, len(ips))
			copy(ipsCp, ips)
			return ipsCp
		}
	}
	return nil
}
func (h *Hostsfile) LookupStaticHostV6(host string) []net.IP {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h.RLock()
	defer h.RUnlock()
	if len(h.hmap.byNameV6) != 0 {
		if ips, ok := h.hmap.byNameV6[absDomainName(host)]; ok {
			ipsCp := make([]net.IP, len(ips))
			copy(ipsCp, ips)
			return ipsCp
		}
	}
	return nil
}
func (h *Hostsfile) LookupStaticAddr(addr string) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h.RLock()
	defer h.RUnlock()
	addr = parseLiteralIP(addr).String()
	if addr == "" {
		return nil
	}
	if len(h.hmap.byAddr) != 0 {
		if hosts, ok := h.hmap.byAddr[addr]; ok {
			hostsCp := make([]string, len(hosts))
			copy(hostsCp, hosts)
			return hostsCp
		}
	}
	return nil
}
