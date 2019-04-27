package file

import (
	"github.com/coredns/coredns/plugin/pkg/dnsutil"
	"github.com/miekg/dns"
)

func substituteDNAME(qname, owner, target string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if dns.IsSubDomain(owner, qname) && qname != owner {
		labels := dns.SplitDomainName(qname)
		labels = append(labels[0:len(labels)-dns.CountLabel(owner)], dns.SplitDomainName(target)...)
		return dnsutil.Join(labels...)
	}
	return ""
}
func synthesizeCNAME(qname string, d *dns.DNAME) *dns.CNAME {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	target := substituteDNAME(qname, d.Header().Name, d.Target)
	if target == "" {
		return nil
	}
	r := new(dns.CNAME)
	r.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: d.Header().Ttl}
	r.Target = target
	return r
}
