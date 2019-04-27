package test

import (
	"testing"
	"github.com/coredns/coredns/plugin/proxy"
	"github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func TestProxyErratic(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	corefile := `example.org:0 {
		erratic {
			drop 2
		}
	}
`
	backend, udp, _, err := CoreDNSServerAndPorts(corefile)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}
	defer backend.Stop()
	p := proxy.NewLookup([]string{udp})
	state := request.Request{W: &test.ResponseWriter{}, Req: new(dns.Msg)}
	p.Lookup(state, "example.org.", dns.TypeA)
}
func TestProxyThreeWay(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	corefileUp1 := `example.org:0 {
		erratic {
			drop 2
		}
	}
`
	up1, err := CoreDNSServer(corefileUp1)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}
	defer up1.Stop()
	corefileUp2 := `example.org:0 {
		whoami
	}
`
	up2, err := CoreDNSServer(corefileUp2)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}
	defer up2.Stop()
	addr1, _ := CoreDNSServerPorts(up1, 0)
	if addr1 == "" {
		t.Fatalf("Could not get UDP listening port")
	}
	addr2, _ := CoreDNSServerPorts(up2, 0)
	if addr2 == "" {
		t.Fatalf("Could not get UDP listening port")
	}
	corefileProxy := `example.org:0 {
		proxy . ` + addr1 + " " + addr2 + ` {
			max_fails 1
		}
	}`
	prx, err := CoreDNSServer(corefileProxy)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}
	defer prx.Stop()
	addr, _ := CoreDNSServerPorts(prx, 0)
	if addr == "" {
		t.Fatalf("Could not get UDP listening port")
	}
	m := new(dns.Msg)
	m.SetQuestion("example.org.", dns.TypeA)
	for i := 0; i < 10; i++ {
		r, err := dns.Exchange(m, addr)
		if err != nil {
			continue
		}
		if r.Rcode != dns.RcodeSuccess {
			t.Fatalf("Expected success rcode, got %d", r.Rcode)
		}
	}
}
