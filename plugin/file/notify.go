package file

import (
	"fmt"
	"net"
	"github.com/coredns/coredns/plugin/pkg/rcode"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func (z *Zone) isNotify(state request.Request) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if state.Req.Opcode != dns.OpcodeNotify {
		return false
	}
	if len(z.TransferFrom) == 0 {
		return false
	}
	remote := state.IP()
	for _, f := range z.TransferFrom {
		from, _, err := net.SplitHostPort(f)
		if err != nil {
			continue
		}
		if from == remote {
			return true
		}
	}
	return false
}
func (z *Zone) Notify() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	go notify(z.origin, z.TransferTo)
}
func notify(zone string, to []string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := new(dns.Msg)
	m.SetNotify(zone)
	c := new(dns.Client)
	for _, t := range to {
		if t == "*" {
			continue
		}
		if err := notifyAddr(c, m, t); err != nil {
			log.Error(err.Error())
		} else {
			log.Infof("Sent notify for zone %q to %q", zone, t)
		}
	}
	return nil
}
func notifyAddr(c *dns.Client, m *dns.Msg, s string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var err error
	code := dns.RcodeServerFailure
	for i := 0; i < 3; i++ {
		ret, _, err := c.Exchange(m, s)
		if err != nil {
			continue
		}
		code = ret.Rcode
		if code == dns.RcodeSuccess {
			return nil
		}
	}
	if err != nil {
		return fmt.Errorf("notify for zone %q was not accepted by %q: %q", m.Question[0].Name, s, err)
	}
	return fmt.Errorf("notify for zone %q was not accepted by %q: rcode was %q", m.Question[0].Name, s, rcode.ToString(code))
}
