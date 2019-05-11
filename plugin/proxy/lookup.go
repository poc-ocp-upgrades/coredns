package proxy

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"time"
	"github.com/coredns/coredns/plugin/pkg/healthcheck"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func NewLookup(hosts []string) Proxy {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return NewLookupWithOption(hosts, Options{})
}
func NewLookupWithOption(hosts []string, opts Options) Proxy {
	_logClusterCodePath()
	defer _logClusterCodePath()
	p := Proxy{Next: nil}
	upstream := &staticUpstream{from: ".", HealthCheck: healthcheck.HealthCheck{FailTimeout: 5 * time.Second, MaxFails: 3}, ex: newDNSExWithOption(opts)}
	upstream.Hosts = make([]*healthcheck.UpstreamHost, len(hosts))
	for i, host := range hosts {
		uh := &healthcheck.UpstreamHost{Name: host, FailTimeout: upstream.FailTimeout, CheckDown: checkDownFunc(upstream)}
		upstream.Hosts[i] = uh
	}
	p.Upstreams = &[]Upstream{upstream}
	return p
}
func (p Proxy) Lookup(state request.Request, name string, typ uint16) (*dns.Msg, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req := new(dns.Msg)
	req.SetQuestion(name, typ)
	state.SizeAndDo(req)
	state2 := request.Request{W: state.W, Req: req}
	return p.lookup(state2)
}
func (p Proxy) Forward(state request.Request) (*dns.Msg, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return p.lookup(state)
}
func (p Proxy) lookup(state request.Request) (*dns.Msg, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	upstream := p.match(state)
	if upstream == nil {
		return nil, errInvalidDomain
	}
	for {
		start := time.Now()
		var reply *dns.Msg
		var backendErr error
		for time.Since(start) < tryDuration {
			host := upstream.Select()
			if host == nil {
				return nil, fmt.Errorf("%s: %s", errUnreachable, "no upstream host")
			}
			atomic.AddInt64(&host.Conns, 1)
			reply, backendErr = upstream.Exchanger().Exchange(context.TODO(), host.Name, state)
			atomic.AddInt64(&host.Conns, -1)
			if backendErr == nil {
				if !state.Match(reply) {
					return state.ErrorMessage(dns.RcodeFormatError), nil
				}
				return reply, nil
			}
			if oe, ok := backendErr.(*net.OpError); ok {
				if oe.Timeout() {
					continue
				}
			}
			timeout := host.FailTimeout
			if timeout == 0 {
				timeout = defaultFailTimeout
			}
			atomic.AddInt32(&host.Fails, 1)
			fails := atomic.LoadInt32(&host.Fails)
			go func(host *healthcheck.UpstreamHost, timeout time.Duration) {
				time.Sleep(timeout)
				atomic.AddInt32(&host.Fails, -1)
				if fails%failureCheck == 0 {
					host.HealthCheckURL()
				}
			}(host, timeout)
		}
		return nil, fmt.Errorf("%s: %s", errUnreachable, backendErr)
	}
}
