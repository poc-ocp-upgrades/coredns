package cache

import (
	"context"
	"testing"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
)

func TestFormErr(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := New()
	c.Next = formErrHandler()
	req := new(dns.Msg)
	req.SetQuestion("example.org.", dns.TypeA)
	rec := dnstest.NewRecorder(&test.ResponseWriter{})
	c.ServeDNS(context.TODO(), rec, req)
	if c.pcache.Len() != 0 {
		t.Errorf("Cached %s, while reply had %d", "example.org.", rec.Msg.Rcode)
	}
}
func formErrHandler() plugin.Handler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return plugin.HandlerFunc(func(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
		m := new(dns.Msg)
		m.SetQuestion("example.net.", dns.TypeA)
		m.Rcode = dns.RcodeFormatError
		w.WriteMsg(m)
		return dns.RcodeSuccess, nil
	})
}
