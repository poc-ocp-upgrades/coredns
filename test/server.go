package test

import (
	"sync"
	"github.com/coredns/coredns/core/dnsserver"
	_ "github.com/coredns/coredns/core"
	"github.com/mholt/caddy"
)

var mu sync.Mutex

func CoreDNSServer(corefile string) (*caddy.Instance, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mu.Lock()
	defer mu.Unlock()
	caddy.Quiet = true
	dnsserver.Quiet = true
	return caddy.Start(NewInput(corefile))
}
func CoreDNSServerStop(i *caddy.Instance) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	i.Stop()
}
func CoreDNSServerPorts(i *caddy.Instance, k int) (udp, tcp string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	srvs := i.Servers()
	if len(srvs) < k+1 {
		return "", ""
	}
	u := srvs[k].LocalAddr()
	t := srvs[k].Addr()
	if u != nil {
		udp = u.String()
	}
	if t != nil {
		tcp = t.String()
	}
	return
}
func CoreDNSServerAndPorts(corefile string) (i *caddy.Instance, udp, tcp string, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	i, err = CoreDNSServer(corefile)
	if err != nil {
		return nil, "", "", err
	}
	udp, tcp = CoreDNSServerPorts(i, 0)
	return i, udp, tcp, nil
}

type Input struct{ corefile []byte }

func NewInput(corefile string) *Input {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &Input{corefile: []byte(corefile)}
}
func (i *Input) Body() []byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return i.corefile
}
func (i *Input) Path() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "Corefile"
}
func (i *Input) ServerType() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "dns"
}
