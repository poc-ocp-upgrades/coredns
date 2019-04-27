package forward

import (
	"crypto/tls"
	"sync/atomic"
	"time"
	"github.com/coredns/coredns/plugin/pkg/transport"
	"github.com/miekg/dns"
)

type HealthChecker interface {
	Check(*Proxy) error
	SetTLSConfig(*tls.Config)
}
type dnsHc struct{ c *dns.Client }

func NewHealthChecker(trans string) HealthChecker {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch trans {
	case transport.DNS, transport.TLS:
		c := new(dns.Client)
		c.Net = "udp"
		c.ReadTimeout = 1 * time.Second
		c.WriteTimeout = 1 * time.Second
		return &dnsHc{c: c}
	}
	log.Warningf("No healthchecker for transport %q", trans)
	return nil
}
func (h *dnsHc) SetTLSConfig(cfg *tls.Config) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	h.c.Net = "tcp-tls"
	h.c.TLSConfig = cfg
}
func (h *dnsHc) Check(p *Proxy) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	err := h.send(p.addr)
	if err != nil {
		HealthcheckFailureCount.WithLabelValues(p.addr).Add(1)
		atomic.AddUint32(&p.fails, 1)
		return err
	}
	atomic.StoreUint32(&p.fails, 0)
	return nil
}
func (h *dnsHc) send(addr string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ping := new(dns.Msg)
	ping.SetQuestion(".", dns.TypeNS)
	m, _, err := h.c.Exchange(ping, addr)
	if err != nil && m != nil {
		if m.Response || m.Opcode == dns.OpcodeQuery {
			err = nil
		}
	}
	return err
}
