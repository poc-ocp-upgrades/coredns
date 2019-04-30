package kubernetes

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func (k Kubernetes) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	opt := plugin.Options{}
	state := request.Request{W: w, Req: r, Context: ctx}
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	qname := state.QName()
	zone := plugin.Zones(k.Zones).Matches(qname)
	if zone == "" {
		return plugin.NextOrFailure(k.Name(), k.Next, ctx, w, r)
	}
	zone = qname[len(qname)-len(zone):]
	state.Zone = zone
	var (
		records	[]dns.RR
		extra	[]dns.RR
		err	error
	)
	switch state.QType() {
	case dns.TypeA:
		records, err = plugin.A(&k, zone, state, nil, opt)
	case dns.TypeAAAA:
		records, err = plugin.AAAA(&k, zone, state, nil, opt)
	case dns.TypeTXT:
		records, err = plugin.TXT(&k, zone, state, opt)
	case dns.TypeCNAME:
		records, err = plugin.CNAME(&k, zone, state, opt)
	case dns.TypePTR:
		records, err = plugin.PTR(&k, zone, state, opt)
	case dns.TypeMX:
		records, extra, err = plugin.MX(&k, zone, state, opt)
	case dns.TypeSRV:
		records, extra, err = plugin.SRV(&k, zone, state, opt)
	case dns.TypeSOA:
		records, err = plugin.SOA(&k, zone, state, opt)
	case dns.TypeNS:
		if state.Name() == zone {
			records, extra, err = plugin.NS(&k, zone, state, opt)
			break
		}
		fallthrough
	case dns.TypeAXFR, dns.TypeIXFR:
		k.Transfer(ctx, state)
	default:
		_, err = plugin.A(&k, zone, state, nil, opt)
	}
	if k.IsNameError(err) {
		if k.Fall.Through(state.Name()) {
			return plugin.NextOrFailure(k.Name(), k.Next, ctx, w, r)
		}
		if !k.APIConn.HasSynced() {
			return plugin.BackendError(&k, zone, dns.RcodeServerFailure, state, nil, opt)
		}
		return plugin.BackendError(&k, zone, dns.RcodeNameError, state, nil, opt)
	}
	if err != nil {
		return dns.RcodeServerFailure, err
	}
	if len(records) == 0 {
		return plugin.BackendError(&k, zone, dns.RcodeSuccess, state, nil, opt)
	}
	m.Answer = append(m.Answer, records...)
	m.Extra = append(m.Extra, extra...)
	w.WriteMsg(m)
	return dns.RcodeSuccess, nil
}
func (k Kubernetes) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "kubernetes"
}
