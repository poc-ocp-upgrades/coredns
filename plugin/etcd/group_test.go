package etcd

import (
	"testing"
	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
)

func TestGroupLookup(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	etc := newEtcdPlugin()
	for _, serv := range servicesGroup {
		set(t, etc, serv.Key, 0, serv)
		defer delete(t, etc, serv.Key)
	}
	for _, tc := range dnsTestCasesGroup {
		m := tc.Msg()
		rec := dnstest.NewRecorder(&test.ResponseWriter{})
		_, err := etc.ServeDNS(ctxt, rec, m)
		if err != nil {
			t.Errorf("Expected no error, got %v\n", err)
			continue
		}
		resp := rec.Msg
		test.SortAndCheck(t, resp, tc)
	}
}

var servicesGroup = []*msg.Service{{Host: "127.0.0.1", Key: "a.dom.skydns.test.", Group: "g1"}, {Host: "127.0.0.2", Key: "b.sub.dom.skydns.test.", Group: "g1"}, {Host: "127.0.0.1", Key: "a.dom2.skydns.test.", Group: "g1"}, {Host: "127.0.0.2", Key: "b.sub.dom2.skydns.test.", Group: ""}, {Host: "127.0.0.1", Key: "a.dom1.skydns.test.", Group: "g1"}, {Host: "127.0.0.2", Key: "b.sub.dom1.skydns.test.", Group: "g2"}}
var dnsTestCasesGroup = []test.Case{{Qname: "dom.skydns.test.", Qtype: dns.TypeA, Answer: []dns.RR{test.A("dom.skydns.test. 300 IN A 127.0.0.1"), test.A("dom.skydns.test. 300 IN A 127.0.0.2")}}, {Qname: "dom2.skydns.test.", Qtype: dns.TypeA, Answer: []dns.RR{test.A("dom2.skydns.test. 300 IN A 127.0.0.1"), test.A("dom2.skydns.test. 300 IN A 127.0.0.2")}}, {Qname: "dom1.skydns.test.", Qtype: dns.TypeA, Answer: []dns.RR{test.A("dom1.skydns.test. 300 IN A 127.0.0.1")}}}
