package external

import (
	"math"
	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func (e *External) a(services []msg.Service, state request.Request) (records []dns.RR) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	dup := make(map[string]struct{})
	for _, s := range services {
		what, ip := s.HostType()
		switch what {
		case dns.TypeCNAME:
		case dns.TypeA:
			if _, ok := dup[s.Host]; !ok {
				dup[s.Host] = struct{}{}
				rr := s.NewA(state.QName(), ip)
				rr.Hdr.Ttl = e.ttl
				records = append(records, rr)
			}
		case dns.TypeAAAA:
		}
	}
	return records
}
func (e *External) aaaa(services []msg.Service, state request.Request) (records []dns.RR) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	dup := make(map[string]struct{})
	for _, s := range services {
		what, ip := s.HostType()
		switch what {
		case dns.TypeCNAME:
		case dns.TypeA:
		case dns.TypeAAAA:
			if _, ok := dup[s.Host]; !ok {
				dup[s.Host] = struct{}{}
				rr := s.NewAAAA(state.QName(), ip)
				rr.Hdr.Ttl = e.ttl
				records = append(records, rr)
			}
		}
	}
	return records
}
func (e *External) srv(services []msg.Service, state request.Request) (records, extra []dns.RR) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	dup := make(map[item]struct{})
	w := make(map[int]int)
	for _, s := range services {
		weight := 100
		if s.Weight != 0 {
			weight = s.Weight
		}
		if _, ok := w[s.Priority]; !ok {
			w[s.Priority] = weight
			continue
		}
		w[s.Priority] += weight
	}
	for _, s := range services {
		if s.Port == -1 {
			continue
		}
		w1 := 100.0 / float64(w[s.Priority])
		if s.Weight == 0 {
			w1 *= 100
		} else {
			w1 *= float64(s.Weight)
		}
		weight := uint16(math.Floor(w1))
		what, ip := s.HostType()
		switch what {
		case dns.TypeCNAME:
		case dns.TypeA, dns.TypeAAAA:
			addr := s.Host
			s.Host = msg.Domain(s.Key)
			srv := s.NewSRV(state.QName(), weight)
			if ok := isDuplicate(dup, srv.Target, "", srv.Port); !ok {
				records = append(records, srv)
			}
			if ok := isDuplicate(dup, srv.Target, addr, 0); !ok {
				hdr := dns.RR_Header{Name: srv.Target, Rrtype: what, Class: dns.ClassINET, Ttl: e.ttl}
				switch what {
				case dns.TypeA:
					extra = append(extra, &dns.A{Hdr: hdr, A: ip})
				case dns.TypeAAAA:
					extra = append(extra, &dns.AAAA{Hdr: hdr, AAAA: ip})
				}
			}
		}
	}
	return records, extra
}

type item struct {
	name	string
	port	uint16
	addr	string
}

func isDuplicate(m map[item]struct{}, name, addr string, port uint16) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if addr != "" {
		_, ok := m[item{name, 0, addr}]
		if !ok {
			m[item{name, 0, addr}] = struct{}{}
		}
		return ok
	}
	_, ok := m[item{name, port, ""}]
	if !ok {
		m[item{name, port, ""}] = struct{}{}
	}
	return ok
}
