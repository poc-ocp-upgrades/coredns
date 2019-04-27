package dnssec

import "github.com/miekg/dns"

func (k *DNSKEY) newRRSIG(signerName string, ttl, incep, expir uint32) *dns.RRSIG {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	sig := new(dns.RRSIG)
	sig.Hdr.Rrtype = dns.TypeRRSIG
	sig.Algorithm = k.K.Algorithm
	sig.KeyTag = k.tag
	sig.SignerName = signerName
	sig.Hdr.Ttl = ttl
	sig.OrigTtl = origTTL
	sig.Inception = incep
	sig.Expiration = expir
	return sig
}

type rrset struct {
	qname	string
	qtype	uint16
}

func rrSets(rrs []dns.RR) map[rrset][]dns.RR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := make(map[rrset][]dns.RR)
	for _, r := range rrs {
		if r.Header().Rrtype == dns.TypeRRSIG || r.Header().Rrtype == dns.TypeOPT {
			continue
		}
		if s, ok := m[rrset{r.Header().Name, r.Header().Rrtype}]; ok {
			s = append(s, r)
			m[rrset{r.Header().Name, r.Header().Rrtype}] = s
			continue
		}
		s := make([]dns.RR, 1, 3)
		s[0] = r
		m[rrset{r.Header().Name, r.Header().Rrtype}] = s
	}
	if len(m) > 0 {
		return m
	}
	return nil
}

const origTTL = 3600
