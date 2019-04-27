package nsid

import (
	"context"
	"encoding/hex"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

type Nsid struct {
	Next	plugin.Handler
	Data	string
}
type ResponseWriter struct {
	dns.ResponseWriter
	Data	string
}

func (n Nsid) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if option := r.IsEdns0(); option != nil {
		for _, o := range option.Option {
			if _, ok := o.(*dns.EDNS0_NSID); ok {
				nw := &ResponseWriter{ResponseWriter: w, Data: n.Data}
				return plugin.NextOrFailure(n.Name(), n.Next, ctx, nw, r)
			}
		}
	}
	return plugin.NextOrFailure(n.Name(), n.Next, ctx, w, r)
}
func (w *ResponseWriter) WriteMsg(res *dns.Msg) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if option := res.IsEdns0(); option != nil {
		for _, o := range option.Option {
			if e, ok := o.(*dns.EDNS0_NSID); ok {
				e.Code = dns.EDNS0NSID
				e.Nsid = hex.EncodeToString([]byte(w.Data))
			}
		}
	}
	returned := w.ResponseWriter.WriteMsg(res)
	return returned
}
func (n Nsid) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "nsid"
}
