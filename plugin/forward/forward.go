package forward

import (
	"context"
	"crypto/tls"
	"errors"
	"time"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/debug"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	ot "github.com/opentracing/opentracing-go"
)

var log = clog.NewWithPlugin("forward")

type Forward struct {
	proxies		[]*Proxy
	p		Policy
	hcInterval	time.Duration
	from		string
	ignored		[]string
	tlsConfig	*tls.Config
	tlsServerName	string
	maxfails	uint32
	expire		time.Duration
	opts		options
	Next		plugin.Handler
}

func New() *Forward {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	f := &Forward{maxfails: 2, tlsConfig: new(tls.Config), expire: defaultExpire, p: new(random), from: ".", hcInterval: hcInterval}
	return f
}
func (f *Forward) SetProxy(p *Proxy) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	f.proxies = append(f.proxies, p)
	p.start(f.hcInterval)
}
func (f *Forward) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(f.proxies)
}
func (f *Forward) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "forward"
}
func (f *Forward) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	state := request.Request{W: w, Req: r}
	if !f.match(state) {
		return plugin.NextOrFailure(f.Name(), f.Next, ctx, w, r)
	}
	fails := 0
	var span, child ot.Span
	var upstreamErr error
	span = ot.SpanFromContext(ctx)
	i := 0
	list := f.List()
	deadline := time.Now().Add(defaultTimeout)
	for time.Now().Before(deadline) {
		if i >= len(list) {
			i = 0
			fails = 0
		}
		proxy := list[i]
		i++
		if proxy.Down(f.maxfails) {
			fails++
			if fails < len(f.proxies) {
				continue
			}
			r := new(random)
			proxy = r.List(f.proxies)[0]
			HealthcheckBrokenCount.Add(1)
		}
		if span != nil {
			child = span.Tracer().StartSpan("connect", ot.ChildOf(span.Context()))
			ctx = ot.ContextWithSpan(ctx, child)
		}
		var (
			ret	*dns.Msg
			err	error
		)
		opts := f.opts
		for {
			ret, err = proxy.Connect(ctx, state, opts)
			if err == nil {
				break
			}
			if err == ErrCachedClosed {
				continue
			}
			if ret != nil && ret.Truncated && !opts.forceTCP && f.opts.preferUDP {
				opts.forceTCP = true
				continue
			}
			break
		}
		if child != nil {
			child.Finish()
		}
		upstreamErr = err
		if err != nil {
			if f.maxfails != 0 {
				proxy.Healthcheck()
			}
			if fails < len(f.proxies) {
				continue
			}
			break
		}
		if !state.Match(ret) {
			debug.Hexdumpf(ret, "Wrong reply for id: %d, %s/%d", state.QName(), state.QType())
			formerr := state.ErrorMessage(dns.RcodeFormatError)
			w.WriteMsg(formerr)
			return 0, nil
		}
		w.WriteMsg(ret)
		return 0, nil
	}
	if upstreamErr != nil {
		return dns.RcodeServerFailure, upstreamErr
	}
	return dns.RcodeServerFailure, ErrNoHealthy
}
func (f *Forward) match(state request.Request) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !plugin.Name(f.from).Matches(state.Name()) || !f.isAllowedDomain(state.Name()) {
		return false
	}
	return true
}
func (f *Forward) isAllowedDomain(name string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if dns.Name(name) == dns.Name(f.from) {
		return true
	}
	for _, ignore := range f.ignored {
		if plugin.Name(ignore).Matches(name) {
			return false
		}
	}
	return true
}
func (f *Forward) ForceTCP() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return f.opts.forceTCP
}
func (f *Forward) PreferUDP() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return f.opts.preferUDP
}
func (f *Forward) List() []*Proxy {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return f.p.List(f.proxies)
}

var (
	ErrNoHealthy	= errors.New("no healthy proxies")
	ErrNoForward	= errors.New("no forwarder defined")
	ErrCachedClosed	= errors.New("cached connection was closed by peer")
)

type policy int

const (
	randomPolicy	policy	= iota
	roundRobinPolicy
	sequentialPolicy
)

type options struct {
	forceTCP	bool
	preferUDP	bool
}

const defaultTimeout = 5 * time.Second
