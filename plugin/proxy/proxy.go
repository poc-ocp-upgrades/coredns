package proxy

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync/atomic"
	"time"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	"github.com/coredns/coredns/plugin/pkg/healthcheck"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	ot "github.com/opentracing/opentracing-go"
)

var (
	errUnreachable		= errors.New("unreachable backend")
	errInvalidProtocol	= errors.New("invalid protocol")
	errInvalidDomain	= errors.New("invalid path for proxy")
)

type Proxy struct {
	Next		plugin.Handler
	Upstreams	*[]Upstream
	Trace		plugin.Handler
}
type Upstream interface {
	From() string
	Select() *healthcheck.UpstreamHost
	IsAllowedDomain(string) bool
	Exchanger() Exchanger
	Stop() error
}

var tryDuration = 16 * time.Second

func (p Proxy) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var span, child ot.Span
	span = ot.SpanFromContext(ctx)
	state := request.Request{W: w, Req: r}
	upstream := p.match(state)
	if upstream == nil {
		return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
	}
	for {
		start := time.Now()
		var reply *dns.Msg
		var backendErr error
		for time.Since(start) < tryDuration {
			host := upstream.Select()
			if host == nil {
				return dns.RcodeServerFailure, fmt.Errorf("%s: %s", errUnreachable, "no upstream host")
			}
			if span != nil {
				child = span.Tracer().StartSpan("exchange", ot.ChildOf(span.Context()))
				ctx = ot.ContextWithSpan(ctx, child)
			}
			atomic.AddInt64(&host.Conns, 1)
			RequestCount.WithLabelValues(metrics.WithServer(ctx), state.Proto(), upstream.Exchanger().Protocol(), familyToString(state.Family()), host.Name).Add(1)
			reply, backendErr = upstream.Exchanger().Exchange(ctx, host.Name, state)
			atomic.AddInt64(&host.Conns, -1)
			if child != nil {
				child.Finish()
			}
			taperr := toDnstap(ctx, host.Name, upstream.Exchanger(), state, reply, start)
			if backendErr == nil {
				if !state.Match(reply) {
					formerr := state.ErrorMessage(dns.RcodeFormatError)
					w.WriteMsg(formerr)
					return 0, taperr
				}
				w.WriteMsg(reply)
				RequestDuration.WithLabelValues(metrics.WithServer(ctx), state.Proto(), upstream.Exchanger().Protocol(), familyToString(state.Family()), host.Name).Observe(time.Since(start).Seconds())
				return 0, taperr
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
		return dns.RcodeServerFailure, fmt.Errorf("%s: %s", errUnreachable, backendErr)
	}
}
func (p Proxy) match(state request.Request) (u Upstream) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if p.Upstreams == nil {
		return nil
	}
	longestMatch := 0
	for _, upstream := range *p.Upstreams {
		from := upstream.From()
		if !plugin.Name(from).Matches(state.Name()) || !upstream.IsAllowedDomain(state.Name()) {
			continue
		}
		if lf := len(from); lf > longestMatch {
			longestMatch = lf
			u = upstream
		}
	}
	return u
}
func (p Proxy) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "proxy"
}

const (
	defaultFailTimeout	= 2 * time.Second
	defaultTimeout		= 5 * time.Second
	failureCheck		= 3
)
