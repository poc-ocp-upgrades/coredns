package forward

import (
	"crypto/tls"
	"runtime"
	"sync/atomic"
	"time"
	"github.com/coredns/coredns/plugin/pkg/up"
)

type Proxy struct {
	fails		uint32
	addr		string
	expire		time.Duration
	transport	*Transport
	probe		*up.Probe
	health		HealthChecker
}

func NewProxy(addr, trans string) *Proxy {
	_logClusterCodePath()
	defer _logClusterCodePath()
	p := &Proxy{addr: addr, fails: 0, probe: up.New(), transport: newTransport(addr)}
	p.health = NewHealthChecker(trans)
	runtime.SetFinalizer(p, (*Proxy).finalizer)
	return p
}
func (p *Proxy) SetTLSConfig(cfg *tls.Config) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	p.transport.SetTLSConfig(cfg)
	p.health.SetTLSConfig(cfg)
}
func (p *Proxy) SetExpire(expire time.Duration) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	p.transport.SetExpire(expire)
}
func (p *Proxy) Healthcheck() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if p.health == nil {
		log.Warning("No healthchecker")
		return
	}
	p.probe.Do(func() error {
		return p.health.Check(p)
	})
}
func (p *Proxy) Down(maxfails uint32) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if maxfails == 0 {
		return false
	}
	fails := atomic.LoadUint32(&p.fails)
	return fails > maxfails
}
func (p *Proxy) close() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	p.probe.Stop()
}
func (p *Proxy) finalizer() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	p.transport.Stop()
}
func (p *Proxy) start(duration time.Duration) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	p.probe.Start(duration)
	p.transport.Start()
}

const (
	maxTimeout	= 2 * time.Second
	minTimeout	= 200 * time.Millisecond
	hcInterval	= 500 * time.Millisecond
)
