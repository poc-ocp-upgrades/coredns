package proxy

import (
	"github.com/coredns/coredns/plugin/pkg/fuzz"
	"github.com/mholt/caddy"
)

func Fuzz(data []byte) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := caddy.NewTestController("dns", "proxy . 8.8.8.8:53")
	up, err := NewStaticUpstreams(&c.Dispenser)
	if err != nil {
		return 0
	}
	p := &Proxy{Upstreams: &up}
	return fuzz.Do(p, data)
}
