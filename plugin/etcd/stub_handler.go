package etcd

import (
	"context"
	"errors"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Stub struct {
	*Etcd
	Zone	string
}

func (s Stub) ServeDNS(ctx context.Context, w dns.ResponseWriter, req *dns.Msg) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if hasStubEdns0(req) {
		log.Warningf("Forwarding cycle detected, refusing msg: %s", req.Question[0].Name)
		return dns.RcodeRefused, errors.New("stub forward cycle")
	}
	req = addStubEdns0(req)
	proxy, ok := (*s.Etcd.Stubmap)[s.Zone]
	if !ok {
		return dns.RcodeServerFailure, nil
	}
	state := request.Request{W: w, Req: req}
	m, e := proxy.Forward(state)
	if e != nil {
		return dns.RcodeServerFailure, e
	}
	w.WriteMsg(m)
	return dns.RcodeSuccess, nil
}
func hasStubEdns0(m *dns.Msg) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	option := m.IsEdns0()
	if option == nil {
		return false
	}
	for _, o := range option.Option {
		if o.Option() == ednsStubCode && len(o.(*dns.EDNS0_LOCAL).Data) == 1 && o.(*dns.EDNS0_LOCAL).Data[0] == 1 {
			return true
		}
	}
	return false
}
func addStubEdns0(m *dns.Msg) *dns.Msg {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	option := m.IsEdns0()
	if option != nil {
		option.Option = append(option.Option, &dns.EDNS0_LOCAL{Code: ednsStubCode, Data: []byte{1}})
		return m
	}
	m.Extra = append(m.Extra, ednsStub)
	return m
}

const (
	ednsStubCode	= dns.EDNS0LOCALSTART + 10
	stubDomain	= "stub.dns"
)

var ednsStub = func() *dns.OPT {
	o := new(dns.OPT)
	o.Hdr.Name = "."
	o.Hdr.Rrtype = dns.TypeOPT
	o.SetUDPSize(4096)
	e := new(dns.EDNS0_LOCAL)
	e.Code = ednsStubCode
	e.Data = []byte{1}
	o.Option = append(o.Option, e)
	return o
}()
