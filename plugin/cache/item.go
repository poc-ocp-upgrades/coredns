package cache

import (
	"time"
	"github.com/coredns/coredns/plugin/cache/freq"
	"github.com/miekg/dns"
)

type item struct {
	Rcode			int
	Authoritative		bool
	AuthenticatedData	bool
	RecursionAvailable	bool
	Answer			[]dns.RR
	Ns			[]dns.RR
	Extra			[]dns.RR
	origTTL			uint32
	stored			time.Time
	*freq.Freq
}

func newItem(m *dns.Msg, now time.Time, d time.Duration) *item {
	_logClusterCodePath()
	defer _logClusterCodePath()
	i := new(item)
	i.Rcode = m.Rcode
	i.Authoritative = m.Authoritative
	i.AuthenticatedData = m.AuthenticatedData
	i.RecursionAvailable = m.RecursionAvailable
	i.Answer = m.Answer
	i.Ns = m.Ns
	i.Extra = make([]dns.RR, len(m.Extra))
	j := 0
	for _, e := range m.Extra {
		if e.Header().Rrtype == dns.TypeOPT {
			continue
		}
		i.Extra[j] = e
		j++
	}
	i.Extra = i.Extra[:j]
	i.origTTL = uint32(d.Seconds())
	i.stored = now.UTC()
	i.Freq = new(freq.Freq)
	return i
}
func (i *item) toMsg(m *dns.Msg, now time.Time) *dns.Msg {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m1 := new(dns.Msg)
	m1.SetReply(m)
	m1.Authoritative = false
	m1.AuthenticatedData = i.AuthenticatedData
	m1.RecursionAvailable = i.RecursionAvailable
	m1.Rcode = i.Rcode
	m1.Answer = make([]dns.RR, len(i.Answer))
	m1.Ns = make([]dns.RR, len(i.Ns))
	m1.Extra = make([]dns.RR, len(i.Extra))
	ttl := uint32(i.ttl(now))
	for j, r := range i.Answer {
		m1.Answer[j] = dns.Copy(r)
		m1.Answer[j].Header().Ttl = ttl
	}
	for j, r := range i.Ns {
		m1.Ns[j] = dns.Copy(r)
		m1.Ns[j].Header().Ttl = ttl
	}
	for j, r := range i.Extra {
		m1.Extra[j] = dns.Copy(r)
		m1.Extra[j].Header().Ttl = ttl
	}
	return m1
}
func (i *item) ttl(now time.Time) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ttl := int(i.origTTL) - int(now.UTC().Sub(i.stored).Seconds())
	return ttl
}
