package dnsutil

import (
	"time"
	"github.com/coredns/coredns/plugin/pkg/response"
	"github.com/miekg/dns"
)

func MinimalTTL(m *dns.Msg, mt response.Type) time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if mt != response.NoError && mt != response.NameError && mt != response.NoData {
		return MinimalDefaultTTL
	}
	if len(m.Answer)+len(m.Ns) == 0 && (len(m.Extra) == 0 || (len(m.Extra) == 1 && m.Extra[0].Header().Rrtype == dns.TypeOPT)) {
		return MinimalDefaultTTL
	}
	minTTL := MaximumDefaulTTL
	for _, r := range m.Answer {
		if r.Header().Ttl < uint32(minTTL.Seconds()) {
			minTTL = time.Duration(r.Header().Ttl) * time.Second
		}
	}
	for _, r := range m.Ns {
		if r.Header().Ttl < uint32(minTTL.Seconds()) {
			minTTL = time.Duration(r.Header().Ttl) * time.Second
		}
	}
	for _, r := range m.Extra {
		if r.Header().Rrtype == dns.TypeOPT {
			continue
		}
		if r.Header().Ttl < uint32(minTTL.Seconds()) {
			minTTL = time.Duration(r.Header().Ttl) * time.Second
		}
	}
	return minTTL
}

const (
	MinimalDefaultTTL	= 5 * time.Second
	MaximumDefaulTTL	= 1 * time.Hour
)
