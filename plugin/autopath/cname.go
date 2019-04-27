package autopath

import (
	"strings"
	"github.com/miekg/dns"
)

func cnamer(m *dns.Msg, original string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, a := range m.Answer {
		if strings.EqualFold(original, a.Header().Name) {
			continue
		}
		m.Answer = append(m.Answer, nil)
		copy(m.Answer[1:], m.Answer)
		m.Answer[0] = &dns.CNAME{Hdr: dns.RR_Header{Name: original, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: a.Header().Ttl}, Target: a.Header().Name}
		break
	}
	m.Question[0].Name = original
}
