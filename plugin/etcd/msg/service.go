package msg

import (
	"net"
	"strings"
	"github.com/miekg/dns"
)

type Service struct {
	Host		string	`json:"host,omitempty"`
	Port		int	`json:"port,omitempty"`
	Priority	int	`json:"priority,omitempty"`
	Weight		int	`json:"weight,omitempty"`
	Text		string	`json:"text,omitempty"`
	Mail		bool	`json:"mail,omitempty"`
	TTL		uint32	`json:"ttl,omitempty"`
	TargetStrip	int	`json:"targetstrip,omitempty"`
	Group		string	`json:"group,omitempty"`
	Key		string	`json:"-"`
}

func (s *Service) NewSRV(name string, weight uint16) *dns.SRV {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	host := dns.Fqdn(s.Host)
	if s.TargetStrip > 0 {
		host = targetStrip(host, s.TargetStrip)
	}
	return &dns.SRV{Hdr: dns.RR_Header{Name: name, Rrtype: dns.TypeSRV, Class: dns.ClassINET, Ttl: s.TTL}, Priority: uint16(s.Priority), Weight: weight, Port: uint16(s.Port), Target: host}
}
func (s *Service) NewMX(name string) *dns.MX {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	host := dns.Fqdn(s.Host)
	if s.TargetStrip > 0 {
		host = targetStrip(host, s.TargetStrip)
	}
	return &dns.MX{Hdr: dns.RR_Header{Name: name, Rrtype: dns.TypeMX, Class: dns.ClassINET, Ttl: s.TTL}, Preference: uint16(s.Priority), Mx: host}
}
func (s *Service) NewA(name string, ip net.IP) *dns.A {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &dns.A{Hdr: dns.RR_Header{Name: name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: s.TTL}, A: ip}
}
func (s *Service) NewAAAA(name string, ip net.IP) *dns.AAAA {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &dns.AAAA{Hdr: dns.RR_Header{Name: name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: s.TTL}, AAAA: ip}
}
func (s *Service) NewCNAME(name string, target string) *dns.CNAME {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &dns.CNAME{Hdr: dns.RR_Header{Name: name, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: s.TTL}, Target: dns.Fqdn(target)}
}
func (s *Service) NewTXT(name string) *dns.TXT {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &dns.TXT{Hdr: dns.RR_Header{Name: name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: s.TTL}, Txt: split255(s.Text)}
}
func (s *Service) NewPTR(name string, target string) *dns.PTR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &dns.PTR{Hdr: dns.RR_Header{Name: name, Rrtype: dns.TypePTR, Class: dns.ClassINET, Ttl: s.TTL}, Ptr: dns.Fqdn(target)}
}
func (s *Service) NewNS(name string) *dns.NS {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	host := dns.Fqdn(s.Host)
	if s.TargetStrip > 0 {
		host = targetStrip(host, s.TargetStrip)
	}
	return &dns.NS{Hdr: dns.RR_Header{Name: name, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: s.TTL}, Ns: host}
}
func Group(sx []Service) []Service {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(sx) == 0 {
		return sx
	}
	group := sx[0].Group
	slashes := strings.Count(sx[0].Key, "/")
	length := make([]int, len(sx))
	for i, s := range sx {
		x := strings.Count(s.Key, "/")
		length[i] = x
		if x < slashes {
			if s.Group == "" {
				break
			}
			slashes = x
			group = s.Group
		}
	}
	if group == "" {
		return sx
	}
	ret := []Service{}
	for i, s := range sx {
		if s.Group == "" {
			ret = append(ret, s)
			continue
		}
		if length[i] == slashes && s.Group != group {
			return sx
		}
		if s.Group == group {
			ret = append(ret, s)
		}
	}
	return ret
}
func split255(s string) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(s) < 255 {
		return []string{s}
	}
	sx := []string{}
	p, i := 0, 255
	for {
		if i <= len(s) {
			sx = append(sx, s[p:i])
		} else {
			sx = append(sx, s[p:])
			break
		}
		p, i = p+255, i+255
	}
	return sx
}
func targetStrip(name string, targetStrip int) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	offset, end := 0, false
	for i := 0; i < targetStrip; i++ {
		offset, end = dns.NextLabel(name, offset)
	}
	if end {
		offset = 0
	}
	name = name[offset:]
	return name
}
