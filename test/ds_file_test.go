package test

import (
	"testing"
	"github.com/coredns/coredns/plugin/proxy"
	mtest "github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

var dsTestCases = []mtest.Case{{Qname: "_udp.miek.nl.", Qtype: dns.TypeDS, Rcode: dns.RcodeNameError, Ns: []dns.RR{mtest.SOA("miek.nl.	1800	IN	SOA	linode.atoom.net. miek.miek.nl. 1282630057 14400 3600 604800 14400")}}, {Qname: "miek.nl.", Qtype: dns.TypeDS, Ns: []dns.RR{mtest.SOA("miek.nl.	1800	IN	SOA	linode.atoom.net. miek.miek.nl. 1282630057 14400 3600 604800 14400")}}}

func TestLookupDS(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.Parallel()
	name, rm, err := TempFile(".", miekNL)
	if err != nil {
		t.Fatalf("Failed to create zone: %s", err)
	}
	defer rm()
	corefile := `miek.nl:0 {
       file ` + name + `
}
`
	i, udp, _, err := CoreDNSServerAndPorts(corefile)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}
	defer i.Stop()
	p := proxy.NewLookup([]string{udp})
	state := request.Request{W: &mtest.ResponseWriter{}, Req: new(dns.Msg)}
	for _, tc := range dsTestCases {
		resp, err := p.Lookup(state, tc.Qname, tc.Qtype)
		if err != nil || resp == nil {
			t.Fatalf("Expected to receive reply, but didn't for %s %d", tc.Qname, tc.Qtype)
		}
		mtest.SortAndCheck(t, resp, tc)
	}
}
