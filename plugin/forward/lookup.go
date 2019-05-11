package forward

import (
	"context"
	"github.com/coredns/coredns/plugin/pkg/transport"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func (f *Forward) Forward(state request.Request) (*dns.Msg, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if f == nil {
		return nil, ErrNoForward
	}
	fails := 0
	var upstreamErr error
	for _, proxy := range f.List() {
		if proxy.Down(f.maxfails) {
			fails++
			if fails < len(f.proxies) {
				continue
			}
			proxy = f.List()[0]
		}
		ret, err := proxy.Connect(context.Background(), state, f.opts)
		upstreamErr = err
		if err != nil {
			if fails < len(f.proxies) {
				continue
			}
			break
		}
		if !state.Match(ret) {
			return state.ErrorMessage(dns.RcodeFormatError), nil
		}
		ret = state.Scrub(ret)
		return ret, err
	}
	if upstreamErr != nil {
		return nil, upstreamErr
	}
	return nil, ErrNoHealthy
}
func (f *Forward) Lookup(state request.Request, name string, typ uint16) (*dns.Msg, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if f == nil {
		return nil, ErrNoForward
	}
	req := new(dns.Msg)
	req.SetQuestion(name, typ)
	state.SizeAndDo(req)
	state2 := request.Request{W: state.W, Req: req}
	return f.Forward(state2)
}
func NewLookup(addr []string) *Forward {
	_logClusterCodePath()
	defer _logClusterCodePath()
	f := New()
	for i := range addr {
		p := NewProxy(addr[i], transport.DNS)
		f.SetProxy(p)
	}
	return f
}
