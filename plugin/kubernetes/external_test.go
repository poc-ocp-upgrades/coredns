package kubernetes

import (
	"testing"
	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/plugin/kubernetes/object"
	"github.com/coredns/coredns/plugin/pkg/watch"
	"github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	api "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var extCases = []struct {
	Qname	string
	Qtype	uint16
	Msg	[]msg.Service
	Rcode	int
}{{Qname: "svc1.testns.example.org.", Rcode: dns.RcodeSuccess, Msg: []msg.Service{msg.Service{Host: "1.2.3.4", Port: 80, TTL: 5, Key: "/c/org/example/testns/svc1"}}}, {Qname: "svc6.testns.example.org.", Rcode: dns.RcodeSuccess, Msg: []msg.Service{msg.Service{Host: "1:2::5", Port: 80, TTL: 5, Key: "/c/org/example/testns/svc1"}}}, {Qname: "*._not-udp-or-tcp.svc1.testns.example.com.", Rcode: dns.RcodeSuccess}, {Qname: "_http._tcp.svc1.testns.example.com.", Rcode: dns.RcodeSuccess, Msg: []msg.Service{msg.Service{Host: "1.2.3.4", Port: 80, TTL: 5, Key: "/c/org/example/testns/svc1"}}}, {Qname: "svc0.testns.example.com.", Rcode: dns.RcodeNameError}, {Qname: "svc0.svc-nons.example.com.", Rcode: dns.RcodeNameError}}

func TestExternal(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	k := New([]string{"cluster.local."})
	k.APIConn = &external{}
	k.Next = test.NextHandler(dns.RcodeSuccess, nil)
	k.Namespaces = map[string]struct{}{"testns": struct{}{}}
	for i, tc := range extCases {
		state := testRequest(tc.Qname)
		svc, rcode := k.External(state)
		if x := tc.Rcode; x != rcode {
			t.Errorf("Test %d, expected rcode %d, got %d\n", i, x, rcode)
		}
		if len(svc) != len(tc.Msg) {
			t.Errorf("Test %d, expected %d for messages, got %d\n", i, len(tc.Msg), len(svc))
		}
		for j, s := range svc {
			if x := tc.Msg[j].Key; x != s.Key {
				t.Errorf("Test %d, expected key %s, got %s\n", i, x, s.Key)
			}
			return
		}
	}
}

type external struct{}

func (external) HasSynced() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return true
}
func (external) Run() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return
}
func (external) Stop() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (external) EpIndexReverse(string) []*object.Endpoints {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (external) SvcIndexReverse(string) []*object.Service {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (external) Modified() int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return 0
}
func (external) SetWatchChan(watch.Chan) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (external) Watch(string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (external) StopWatching(string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (external) EpIndex(s string) []*object.Endpoints {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (external) EndpointsList() []*object.Endpoints {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (external) GetNodeByName(name string) (*api.Node, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, nil
}
func (external) SvcIndex(s string) []*object.Service {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return svcIndexExternal[s]
}
func (external) PodIndex(string) []*object.Pod {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (external) GetNamespaceByName(name string) (*api.Namespace, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &api.Namespace{ObjectMeta: meta.ObjectMeta{Name: name}}, nil
}

var svcIndexExternal = map[string][]*object.Service{"svc1.testns": {{Name: "svc1", Namespace: "testns", Type: api.ServiceTypeClusterIP, ClusterIP: "10.0.0.1", ExternalIPs: []string{"1.2.3.4"}, Ports: []api.ServicePort{{Name: "http", Protocol: "tcp", Port: 80}}}}, "svc6.testns": {{Name: "svc6", Namespace: "testns", Type: api.ServiceTypeClusterIP, ClusterIP: "10.0.0.3", ExternalIPs: []string{"1:2::5"}, Ports: []api.ServicePort{{Name: "http", Protocol: "tcp", Port: 80}}}}}

func (external) ServiceList() []*object.Service {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var svcs []*object.Service
	for _, svc := range svcIndexExternal {
		svcs = append(svcs, svc...)
	}
	return svcs
}
func testRequest(name string) request.Request {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := new(dns.Msg).SetQuestion(name, dns.TypeA)
	return request.Request{W: &test.ResponseWriter{}, Req: m, Zone: "example.org."}
}
