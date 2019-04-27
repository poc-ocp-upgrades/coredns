package plugin

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"github.com/coredns/coredns/plugin/pkg/parse"
	"github.com/miekg/dns"
)

type Zones []string

func (z Zones) Matches(qname string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	zone := ""
	for _, zname := range z {
		if dns.IsSubDomain(zname, qname) {
			if len(zname) > len(zone) {
				zone = zname
			}
		}
	}
	return zone
}
func (z Zones) Normalize() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for i := range z {
		z[i] = Name(z[i]).Normalize()
	}
}

type Name string

func (n Name) Matches(child string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if dns.Name(n) == dns.Name(child) {
		return true
	}
	return dns.IsSubDomain(string(n), child)
}
func (n Name) Normalize() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return strings.ToLower(dns.Fqdn(string(n)))
}

type (
	Host string
)

func (h Host) Normalize() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := string(h)
	_, s = parse.Transport(s)
	host, _, _, _ := SplitHostPort(s)
	return Name(host).Normalize()
}
func SplitHostPort(s string) (host, port string, ipnet *net.IPNet, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	host = s
	colon := strings.LastIndex(s, ":")
	if colon == len(s)-1 {
		return "", "", nil, fmt.Errorf("expecting data after last colon: %q", s)
	}
	if colon != -1 {
		if p, err := strconv.Atoi(s[colon+1:]); err == nil {
			port = strconv.Itoa(p)
			host = s[:colon]
		}
	}
	if len(host) > 255 {
		return "", "", nil, fmt.Errorf("specified zone is too long: %d > 255", len(host))
	}
	_, d := dns.IsDomainName(host)
	if !d {
		return "", "", nil, fmt.Errorf("zone is not a valid domain name: %s", host)
	}
	ip, n, err := net.ParseCIDR(host)
	ones, bits := 0, 0
	if err == nil {
		if rev, e := dns.ReverseAddr(ip.String()); e == nil {
			ones, bits = n.Mask.Size()
			sizeDigit := 8
			if len(n.IP) == net.IPv6len {
				sizeDigit = 4
			}
			mod := (bits - ones) % sizeDigit
			nearest := (bits - ones) + mod
			offset := 0
			var end bool
			for i := 0; i < nearest/sizeDigit; i++ {
				offset, end = dns.NextLabel(rev, offset)
				if end {
					break
				}
			}
			host = rev[offset:]
		}
	}
	return host, port, n, nil
}
