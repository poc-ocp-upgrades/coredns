package loadbalance

import (
	"github.com/miekg/dns"
)

type RoundRobinResponseWriter struct{ dns.ResponseWriter }

func (r *RoundRobinResponseWriter) WriteMsg(res *dns.Msg) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if res.Rcode != dns.RcodeSuccess {
		return r.ResponseWriter.WriteMsg(res)
	}
	if res.Question[0].Qtype == dns.TypeAXFR || res.Question[0].Qtype == dns.TypeIXFR {
		return r.ResponseWriter.WriteMsg(res)
	}
	res.Answer = roundRobin(res.Answer)
	res.Ns = roundRobin(res.Ns)
	res.Extra = roundRobin(res.Extra)
	return r.ResponseWriter.WriteMsg(res)
}
func roundRobin(in []dns.RR) []dns.RR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cname := []dns.RR{}
	address := []dns.RR{}
	mx := []dns.RR{}
	rest := []dns.RR{}
	for _, r := range in {
		switch r.Header().Rrtype {
		case dns.TypeCNAME:
			cname = append(cname, r)
		case dns.TypeA, dns.TypeAAAA:
			address = append(address, r)
		case dns.TypeMX:
			mx = append(mx, r)
		default:
			rest = append(rest, r)
		}
	}
	roundRobinShuffle(address)
	roundRobinShuffle(mx)
	out := append(cname, rest...)
	out = append(out, address...)
	out = append(out, mx...)
	return out
}
func roundRobinShuffle(records []dns.RR) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch l := len(records); l {
	case 0, 1:
		break
	case 2:
		if dns.Id()%2 == 0 {
			records[0], records[1] = records[1], records[0]
		}
	default:
		for j := 0; j < l*(int(dns.Id())%4+1); j++ {
			q := int(dns.Id()) % l
			p := int(dns.Id()) % l
			if q == p {
				p = (p + 1) % l
			}
			records[q], records[p] = records[p], records[q]
		}
	}
}
func (r *RoundRobinResponseWriter) Write(buf []byte) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	log.Warning("RoundRobin called with Write: not shuffling records")
	n, err := r.ResponseWriter.Write(buf)
	return n, err
}
