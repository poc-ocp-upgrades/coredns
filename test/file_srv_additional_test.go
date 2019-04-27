package test

import (
	"testing"
	"github.com/coredns/coredns/plugin/proxy"
	"github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func TestZoneSRVAdditional(t *testing.T) {
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
	resp, err := p.Lookup(state, "service.example.org.", dns.TypeSRV)
	if err != nil {
		t.Fatalf("Expected to receive reply, but didn't: %s", err)
	}
	if len(resp.Extra) != 2 {
		t.Fatalf("Expected 2 RR in additional section got %d", len(resp.Extra))
	}
}
