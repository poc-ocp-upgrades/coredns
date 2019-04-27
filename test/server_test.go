package test

import (
	"testing"
	"github.com/miekg/dns"
)

func TestProxyToChaosServer(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.Parallel()
	corefile := `.:0 {
	chaos CoreDNS-001 miek@miek.nl
}
`
	chaos, udpChaos, _, err := CoreDNSServerAndPorts(corefile)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}
	defer chaos.Stop()
	corefileProxy := `.:0 {
		proxy . ` + udpChaos + `
}
`
	proxy, udp, _, err := CoreDNSServerAndPorts(corefileProxy)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance")
	}
	defer proxy.Stop()
	chaosTest(t, udpChaos)
	chaosTest(t, udp)
}
func chaosTest(t *testing.T, server string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := new(dns.Msg)
	m.Question = make([]dns.Question, 1)
	m.Question[0] = dns.Question{Qclass: dns.ClassCHAOS, Name: "version.bind.", Qtype: dns.TypeTXT}
	r, err := dns.Exchange(m, server)
	if err != nil {
		t.Fatalf("Could not send message: %s", err)
	}
	if r.Rcode != dns.RcodeSuccess || len(r.Answer) == 0 {
		t.Fatalf("Expected successful reply, got %s", dns.RcodeToString[r.Rcode])
	}
	if r.Answer[0].String() != `version.bind.	0	CH	TXT	"CoreDNS-001"` {
		t.Fatalf("Expected version.bind. reply, got %s", r.Answer[0].String())
	}
}
