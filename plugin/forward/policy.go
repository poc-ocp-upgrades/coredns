package forward

import (
	"math/rand"
	"sync/atomic"
)

type Policy interface {
	List([]*Proxy) []*Proxy
	String() string
}
type random struct{}

func (r *random) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "random"
}
func (r *random) List(p []*Proxy) []*Proxy {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch len(p) {
	case 1:
		return p
	case 2:
		if rand.Int()%2 == 0 {
			return []*Proxy{p[1], p[0]}
		}
		return p
	}
	perms := rand.Perm(len(p))
	rnd := make([]*Proxy, len(p))
	for i, p1 := range perms {
		rnd[i] = p[p1]
	}
	return rnd
}

type roundRobin struct{ robin uint32 }

func (r *roundRobin) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "round_robin"
}
func (r *roundRobin) List(p []*Proxy) []*Proxy {
	_logClusterCodePath()
	defer _logClusterCodePath()
	poolLen := uint32(len(p))
	i := atomic.AddUint32(&r.robin, 1) % poolLen
	robin := []*Proxy{p[i]}
	robin = append(robin, p[:i]...)
	robin = append(robin, p[i+1:]...)
	return robin
}

type sequential struct{}

func (r *sequential) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "sequential"
}
func (r *sequential) List(p []*Proxy) []*Proxy {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return p
}
