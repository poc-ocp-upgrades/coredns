package etcd

import (
	"net"
	"strconv"
	"testing"
	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
)

func fakeStubServerExampleNet(t *testing.T) (*dns.Server, string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	server, addr, err := test.UDPServer("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create a UDP server: %s", err)
	}
	dns.HandleFunc("example.net.", func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Answer = []dns.RR{test.A("example.net.	86400	IN	A	93.184.216.34")}
		w.WriteMsg(m)
	})
	return server, addr
}
func TestStubLookup(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	server, addr := fakeStubServerExampleNet(t)
	defer server.Shutdown()
	host, p, _ := net.SplitHostPort(addr)
	port, _ := strconv.Atoi(p)
	exampleNetStub := &msg.Service{Host: host, Port: port, Key: "a.example.net.stub.dns.skydns.test."}
	servicesStub = append(servicesStub, exampleNetStub)
	etc := newEtcdPlugin()
	for _, serv := range servicesStub {
		set(t, etc, serv.Key, 0, serv)
		defer delete(t, etc, serv.Key)
	}
	etc.updateStubZones()
	for _, tc := range dnsTestCasesStub {
		m := tc.Msg()
		rec := dnstest.NewRecorder(&test.ResponseWriter{})
		_, err := etc.ServeDNS(ctxt, rec, m)
		if err != nil && m.Question[0].Name == "example.org." {
			continue
		}
		if err != nil {
			t.Errorf("Expected no error, got %v for %s\n", err, m.Question[0].Name)
		}
		resp := rec.Msg
		if resp == nil {
			continue
		}
		test.SortAndCheck(t, resp, tc)
	}
}

var servicesStub = []*msg.Service{{Host: "127.0.0.1", Port: 666, Key: "b.example.org.stub.dns.skydns.test."}}
var dnsTestCasesStub = []test.Case{{Qname: "example.org.", Qtype: dns.TypeA, Rcode: dns.RcodeServerFailure}, {Qname: "example.net.", Qtype: dns.TypeA, Answer: []dns.RR{test.A("example.net.	86400	IN	A	93.184.216.34")}}}
