package kubernetes

import (
	"testing"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func TestParseRequest(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tests := []struct {
		query		string
		expected	string
	}{{"_http._tcp.webs.mynamespace.svc.inter.webs.tests.", "http.tcp..webs.mynamespace.svc"}, {"*.any.*.any.svc.inter.webs.tests.", "*.any..*.any.svc"}, {"1-2-3-4.webs.mynamespace.svc.inter.webs.tests.", "*.*.1-2-3-4.webs.mynamespace.svc"}, {"inter.webs.tests.", "....."}, {"svc.inter.webs.tests.", "....."}, {"pod.inter.webs.tests.", "....."}}
	for i, tc := range tests {
		m := new(dns.Msg)
		m.SetQuestion(tc.query, dns.TypeA)
		state := request.Request{Zone: zone, Req: m}
		r, e := parseRequest(state)
		if e != nil {
			t.Errorf("Test %d, expected no error, got '%v'.", i, e)
		}
		rs := r.String()
		if rs != tc.expected {
			t.Errorf("Test %d, expected (stringyfied) recordRequest: %s, got %s", i, tc.expected, rs)
		}
	}
}
func TestParseInvalidRequest(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	invalid := []string{"webs.mynamespace.pood.inter.webs.test.", "too.long.for.what.I.am.trying.to.pod.inter.webs.tests."}
	for i, query := range invalid {
		m := new(dns.Msg)
		m.SetQuestion(query, dns.TypeA)
		state := request.Request{Zone: zone, Req: m}
		if _, e := parseRequest(state); e == nil {
			t.Errorf("Test %d: expected error from %s, got none", i, query)
		}
	}
}

const zone = "inter.webs.tests."
