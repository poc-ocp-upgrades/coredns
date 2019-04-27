package erratic

import (
	"strings"
	"github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func allRecords(name string) []dns.RR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var rrs = []dns.RR{test.SOA("xx.		0	IN	SOA	sns.dns.icann.org. noc.dns.icann.org. 2018050825 7200 3600 1209600 3600"), test.NS("xx.		0	IN	NS	b.xx."), test.NS("xx.		0	IN	NS	a.xx."), test.AAAA("a.xx.	0	IN	AAAA	2001:bd8::53"), test.AAAA("b.xx.	0	IN	AAAA	2001:500::54")}
	for _, r := range rrs {
		r.Header().Name = strings.Replace(r.Header().Name, "xx.", name, 1)
		if n, ok := r.(*dns.NS); ok {
			n.Ns = strings.Replace(n.Ns, "xx.", name, 1)
		}
	}
	return rrs
}
func xfr(state request.Request, truncate bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	rrs := allRecords(state.QName())
	ch := make(chan *dns.Envelope)
	tr := new(dns.Transfer)
	go func() {
		if !truncate {
			rrs = append(rrs, rrs[0])
		}
		ch <- &dns.Envelope{RR: rrs}
		close(ch)
	}()
	tr.Out(state.W, state.Req, ch)
	state.W.Hijack()
	return
}
