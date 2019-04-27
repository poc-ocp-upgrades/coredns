package test

import (
	"net"
	"testing"
	"github.com/miekg/dns"
)

func TestCompressScrub(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	corefile := `example.org:0 {
                       erratic {
			  drop 0
			  delay 0
			  large
		       }
		     }`
	i, udp, _, err := CoreDNSServerAndPorts(corefile)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}
	defer i.Stop()
	c, err := net.Dial("udp", udp)
	if err != nil {
		t.Fatalf("Could not dial %s", err)
	}
	m := new(dns.Msg)
	m.SetQuestion("example.org.", dns.TypeA)
	q, _ := m.Pack()
	c.Write(q)
	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	if err != nil || n == 0 {
		t.Errorf("Expected reply, got: %s", err)
		return
	}
	if n >= 512 {
		t.Fatalf("Expected returned packet to be < 512, got %d", n)
	}
	buf = buf[:n]
	found := 0
	for i := 0; i < len(buf)-1; i++ {
		if buf[i]&0xC0 == 0xC0 {
			off := (int(buf[i])^0xC0)<<8 | int(buf[i+1])
			if off == 12 {
				found++
			}
		}
	}
	if found != 30 {
		t.Errorf("Failed to find all compression pointers in the packet, wanted 30, got %d", found)
	}
}
