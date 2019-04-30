package healthcheck

import (
	"math/rand"
	"sync/atomic"
	"github.com/coredns/coredns/plugin/pkg/log"
)

var (
	SupportedPolicies = make(map[string]func() Policy)
)

func RegisterPolicy(name string, policy func() Policy) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	SupportedPolicies[name] = policy
}

type Policy interface {
	Select(pool HostPool) *UpstreamHost
}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	RegisterPolicy("random", func() Policy {
		return &Random{}
	})
	RegisterPolicy("least_conn", func() Policy {
		return &LeastConn{}
	})
	RegisterPolicy("round_robin", func() Policy {
		return &RoundRobin{}
	})
	RegisterPolicy("first", func() Policy {
		return &First{}
	})
	RegisterPolicy("sequential", func() Policy {
		return &First{}
	})
}

type Random struct{}

func (r *Random) Select(pool HostPool) *UpstreamHost {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var randHost *UpstreamHost
	count := 0
	for _, host := range pool {
		if host.Down() {
			continue
		}
		count++
		if count == 1 {
			randHost = host
		} else {
			r := rand.Int() % count
			if r == (count - 1) {
				randHost = host
			}
		}
	}
	return randHost
}

type Spray struct{}

func (r *Spray) Select(pool HostPool) *UpstreamHost {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rnd := rand.Int() % len(pool)
	randHost := pool[rnd]
	log.Warningf("All hosts reported as down, spraying to target: %s", randHost.Name)
	return randHost
}

type LeastConn struct{}

func (r *LeastConn) Select(pool HostPool) *UpstreamHost {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var bestHost *UpstreamHost
	count := 0
	leastConn := int64(1<<63 - 1)
	for _, host := range pool {
		if host.Down() {
			continue
		}
		hostConns := host.Conns
		if hostConns < leastConn {
			bestHost = host
			leastConn = hostConns
			count = 1
		} else if hostConns == leastConn {
			count++
			if count == 1 {
				bestHost = host
			} else {
				r := rand.Int() % count
				if r == (count - 1) {
					bestHost = host
				}
			}
		}
	}
	return bestHost
}

type RoundRobin struct{ Robin uint32 }

func (r *RoundRobin) Select(pool HostPool) *UpstreamHost {
	_logClusterCodePath()
	defer _logClusterCodePath()
	poolLen := uint32(len(pool))
	selection := atomic.AddUint32(&r.Robin, 1) % poolLen
	host := pool[selection]
	for i := uint32(1); host.Down() && i < poolLen; i++ {
		host = pool[(selection+i)%poolLen]
	}
	return host
}

type First struct{}

func (r *First) Select(pool HostPool) *UpstreamHost {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for i := 0; i < len(pool); i++ {
		host := pool[i]
		if host.Down() {
			continue
		}
		return host
	}
	return nil
}
