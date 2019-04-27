package test

import (
	"testing"
	"github.com/coredns/coredns/plugin/proxy"
	"github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func TestZoneExternalCNAMELookupWithoutProxy(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.Parallel()
	name, rm, err := TempFile(".", exampleOrg)
	if err != nil {
		t.Fatalf("Failed to create zone: %s", err)
	}
	defer rm()
	corefile := `example.org:0 {
       file ` + name + `
}
`
	i, udp, _, err := CoreDNSServerAndPorts(corefile)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}
	defer i.Stop()
	p := proxy.NewLookup([]string{udp})
	state := request.Request{W: &test.ResponseWriter{}, Req: new(dns.Msg)}
	resp, err := p.Lookup(state, "cname.example.org.", dns.TypeA)
	if err != nil {
		t.Fatalf("Expected to receive reply, but didn't: %s", err)
	}
	if len(resp.Answer) != 1 {
		t.Fatalf("Expected 1 RR in answer section got %d", len(resp.Answer))
	}
}
func TestZoneExternalCNAMELookupWithProxy(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.Parallel()
	name, rm, err := TempFile(".", exampleOrg)
	if err != nil {
		t.Fatalf("Failed to create zone: %s", err)
	}
	defer rm()
	corefile := `example.org:0 {
       file ` + name + ` {
	       upstream 8.8.8.8
	}
}
`
	i, udp, _, err := CoreDNSServerAndPorts(corefile)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}
	defer i.Stop()
	p := proxy.NewLookup([]string{udp})
	state := request.Request{W: &test.ResponseWriter{}, Req: new(dns.Msg)}
	resp, err := p.Lookup(state, "cname.example.org.", dns.TypeA)
	if err != nil {
		t.Fatalf("Expected to receive reply, but didn't: %s", err)
	}
	if len(resp.Answer) != 2 {
		t.Fatalf("Expected 2 RRs in answer section got %d", len(resp.Answer))
	}
}
