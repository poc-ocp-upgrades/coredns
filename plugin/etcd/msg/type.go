package msg

import (
	"net"
	"github.com/miekg/dns"
)

func (s *Service) HostType() (what uint16, normalized net.IP) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ip := net.ParseIP(s.Host)
	switch {
	case ip == nil:
		return dns.TypeCNAME, nil
	case ip.To4() != nil:
		return dns.TypeA, ip.To4()
	case ip.To4() == nil:
		return dns.TypeAAAA, ip.To16()
	}
	return dns.TypeNone, nil
}
