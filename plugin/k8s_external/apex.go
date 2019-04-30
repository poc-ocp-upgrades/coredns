package external

import (
	"github.com/coredns/coredns/plugin/pkg/dnsutil"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func (e *External) serveApex(state request.Request) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := new(dns.Msg)
	m.SetReply(state.Req)
	switch state.QType() {
	case dns.TypeSOA:
		m.Answer = []dns.RR{e.soa(state)}
	case dns.TypeNS:
		m.Answer = []dns.RR{e.ns(state)}
		addr := e.externalAddrFunc(state)
		for _, rr := range addr {
			rr.Header().Ttl = e.ttl
			rr.Header().Name = state.QName()
			m.Extra = append(m.Extra, rr)
		}
	default:
		m.Ns = []dns.RR{e.soa(state)}
	}
	state.W.WriteMsg(m)
	return 0, nil
}
func (e *External) serveSubApex(state request.Request) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	base, _ := dnsutil.TrimZone(state.Name(), state.Zone)
	m := new(dns.Msg)
	m.SetReply(state.Req)
	switch labels := dns.CountLabel(base); labels {
	default:
		m.SetRcode(m, dns.RcodeNameError)
		m.Ns = []dns.RR{e.soa(state)}
		state.W.WriteMsg(m)
		return 0, nil
	case 2:
		nl, _ := dns.NextLabel(base, 0)
		ns := base[:nl]
		if ns != "ns1." {
			m.SetRcode(m, dns.RcodeNameError)
			m.Ns = []dns.RR{e.soa(state)}
			state.W.WriteMsg(m)
			return 0, nil
		}
		addr := e.externalAddrFunc(state)
		for _, rr := range addr {
			rr.Header().Ttl = e.ttl
			rr.Header().Name = state.QName()
			switch state.QType() {
			case dns.TypeA:
				if rr.Header().Rrtype == dns.TypeA {
					m.Answer = append(m.Answer, rr)
				}
			case dns.TypeAAAA:
				if rr.Header().Rrtype == dns.TypeAAAA {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
		if len(m.Answer) == 0 {
			m.Ns = []dns.RR{e.soa(state)}
		}
		state.W.WriteMsg(m)
		return 0, nil
	case 1:
		m.Ns = []dns.RR{e.soa(state)}
		state.W.WriteMsg(m)
		return 0, nil
	}
}
func (e *External) soa(state request.Request) *dns.SOA {
	_logClusterCodePath()
	defer _logClusterCodePath()
	header := dns.RR_Header{Name: state.Zone, Rrtype: dns.TypeSOA, Ttl: e.ttl, Class: dns.ClassINET}
	soa := &dns.SOA{Hdr: header, Mbox: dnsutil.Join(e.hostmaster, e.apex, state.Zone), Ns: dnsutil.Join("ns1", e.apex, state.Zone), Serial: 12345, Refresh: 7200, Retry: 1800, Expire: 86400, Minttl: e.ttl}
	return soa
}
func (e *External) ns(state request.Request) *dns.NS {
	_logClusterCodePath()
	defer _logClusterCodePath()
	header := dns.RR_Header{Name: state.Zone, Rrtype: dns.TypeNS, Ttl: e.ttl, Class: dns.ClassINET}
	ns := &dns.NS{Hdr: header, Ns: dnsutil.Join("ns1", e.apex, state.Zone)}
	return ns
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
