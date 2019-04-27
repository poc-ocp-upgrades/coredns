package auto

import (
	"sync"
	"github.com/coredns/coredns/plugin/file"
)

type Zones struct {
	Z	map[string]*file.Zone
	names	[]string
	origins	[]string
	sync.RWMutex
}

func (z *Zones) Names() []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	z.RLock()
	n := z.names
	z.RUnlock()
	return n
}
func (z *Zones) Origins() []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return z.origins
}
func (z *Zones) Zones(name string) *file.Zone {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	z.RLock()
	zo := z.Z[name]
	z.RUnlock()
	return zo
}
func (z *Zones) Add(zo *file.Zone, name string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	z.Lock()
	if z.Z == nil {
		z.Z = make(map[string]*file.Zone)
	}
	z.Z[name] = zo
	z.names = append(z.names, name)
	zo.Reload()
	z.Unlock()
}
func (z *Zones) Remove(name string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	z.Lock()
	if zo, ok := z.Z[name]; ok {
		zo.OnShutdown()
	}
	delete(z.Z, name)
	z.names = []string{}
	for n := range z.Z {
		z.names = append(z.names, n)
	}
	z.Unlock()
}
