package external

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Externaler interface {
	External(request.Request) ([]msg.Service, int)
	ExternalAddress(state request.Request) []dns.RR
}
type External struct {
	Next			plugin.Handler
	Zones			[]string
	hostmaster		string
	apex			string
	ttl			uint32
	externalFunc		func(request.Request) ([]msg.Service, int)
	externalAddrFunc	func(request.Request) []dns.RR
}

func New() *External {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	e := &External{hostmaster: "hostmaster", ttl: 5, apex: "dns"}
	return e
}
func (e *External) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	state := request.Request{W: w, Req: r}
	zone := plugin.Zones(e.Zones).Matches(state.Name())
	if zone == "" {
		return plugin.NextOrFailure(e.Name(), e.Next, ctx, w, r)
	}
	if e.externalFunc == nil {
		return plugin.NextOrFailure(e.Name(), e.Next, ctx, w, r)
	}
	state.Zone = zone
	for _, z := range e.Zones {
		if state.Name() == z {
			ret, err := e.serveApex(state)
			return ret, err
		}
		if dns.IsSubDomain(e.apex+"."+z, state.Name()) {
			ret, err := e.serveSubApex(state)
			return ret, err
		}
	}
	svc, rcode := e.externalFunc(state)
	m := new(dns.Msg)
	m.SetReply(state.Req)
	if len(svc) == 0 {
		m.Rcode = rcode
		m.Ns = []dns.RR{e.soa(state)}
		w.WriteMsg(m)
		return 0, nil
	}
	switch state.QType() {
	case dns.TypeA:
		m.Answer = e.a(svc, state)
	case dns.TypeAAAA:
		m.Answer = e.aaaa(svc, state)
	case dns.TypeSRV:
		m.Answer, m.Extra = e.srv(svc, state)
	default:
		m.Ns = []dns.RR{e.soa(state)}
	}
	if len(m.Answer) == 0 {
		m.Ns = []dns.RR{e.soa(state)}
	}
	w.WriteMsg(m)
	return 0, nil
}
func (e *External) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "k8s_external"
}
