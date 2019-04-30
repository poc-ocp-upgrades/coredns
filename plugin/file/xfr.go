package file

import (
	"context"
	"fmt"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Xfr struct{ *Zone }

func (x Xfr) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	state := request.Request{W: w, Req: r}
	if !x.TransferAllowed(state) {
		return dns.RcodeServerFailure, nil
	}
	if state.QType() != dns.TypeAXFR && state.QType() != dns.TypeIXFR {
		return 0, plugin.Error(x.Name(), fmt.Errorf("xfr called with non transfer type: %d", state.QType()))
	}
	records := x.All()
	if len(records) == 0 {
		return dns.RcodeServerFailure, nil
	}
	ch := make(chan *dns.Envelope)
	defer close(ch)
	tr := new(dns.Transfer)
	go tr.Out(w, r, ch)
	j, l := 0, 0
	records = append(records, records[0])
	log.Infof("Outgoing transfer of %d records of zone %s to %s started", len(records), x.origin, state.IP())
	for i, r := range records {
		l += dns.Len(r)
		if l > transferLength {
			ch <- &dns.Envelope{RR: records[j:i]}
			l = 0
			j = i
		}
	}
	if j < len(records) {
		ch <- &dns.Envelope{RR: records[j:]}
	}
	w.Hijack()
	return dns.RcodeSuccess, nil
}
func (x Xfr) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "xfr"
}

const transferLength = 1000
