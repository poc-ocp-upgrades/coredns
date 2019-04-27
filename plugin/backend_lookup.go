package plugin

import (
	"fmt"
	"math"
	"net"
	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/plugin/pkg/dnsutil"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func A(b ServiceBackend, zone string, state request.Request, previousRecords []dns.RR, opt Options) (records []dns.RR, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	services, err := checkForApex(b, zone, state, opt)
	if err != nil {
		return nil, err
	}
	dup := make(map[string]struct{})
	for _, serv := range services {
		what, ip := serv.HostType()
		switch what {
		case dns.TypeCNAME:
			if Name(state.Name()).Matches(dns.Fqdn(serv.Host)) {
				continue
			}
			newRecord := serv.NewCNAME(state.QName(), serv.Host)
			if len(previousRecords) > 7 {
				continue
			}
			if dnsutil.DuplicateCNAME(newRecord, previousRecords) {
				continue
			}
			if dns.IsSubDomain(zone, dns.Fqdn(serv.Host)) {
				state1 := state.NewWithQuestion(serv.Host, state.QType())
				state1.Zone = zone
				nextRecords, err := A(b, zone, state1, append(previousRecords, newRecord), opt)
				if err == nil {
					if len(nextRecords) > 0 {
						records = append(records, newRecord)
						records = append(records, nextRecords...)
					}
				}
				continue
			}
			target := newRecord.Target
			m1, e1 := b.Lookup(state, target, state.QType())
			if e1 != nil {
				continue
			}
			records = append(records, newRecord)
			records = append(records, m1.Answer...)
			continue
		case dns.TypeA:
			if _, ok := dup[serv.Host]; !ok {
				dup[serv.Host] = struct{}{}
				records = append(records, serv.NewA(state.QName(), ip))
			}
		case dns.TypeAAAA:
		}
	}
	return records, nil
}
func AAAA(b ServiceBackend, zone string, state request.Request, previousRecords []dns.RR, opt Options) (records []dns.RR, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	services, err := checkForApex(b, zone, state, opt)
	if err != nil {
		return nil, err
	}
	dup := make(map[string]struct{})
	for _, serv := range services {
		what, ip := serv.HostType()
		switch what {
		case dns.TypeCNAME:
			if Name(state.Name()).Matches(dns.Fqdn(serv.Host)) {
				continue
			}
			newRecord := serv.NewCNAME(state.QName(), serv.Host)
			if len(previousRecords) > 7 {
				continue
			}
			if dnsutil.DuplicateCNAME(newRecord, previousRecords) {
				continue
			}
			if dns.IsSubDomain(zone, dns.Fqdn(serv.Host)) {
				state1 := state.NewWithQuestion(serv.Host, state.QType())
				state1.Zone = zone
				nextRecords, err := AAAA(b, zone, state1, append(previousRecords, newRecord), opt)
				if err == nil {
					if len(nextRecords) > 0 {
						records = append(records, newRecord)
						records = append(records, nextRecords...)
					}
				}
				continue
			}
			target := newRecord.Target
			m1, e1 := b.Lookup(state, target, state.QType())
			if e1 != nil {
				continue
			}
			records = append(records, newRecord)
			records = append(records, m1.Answer...)
			continue
		case dns.TypeA:
		case dns.TypeAAAA:
			if _, ok := dup[serv.Host]; !ok {
				dup[serv.Host] = struct{}{}
				records = append(records, serv.NewAAAA(state.QName(), ip))
			}
		}
	}
	return records, nil
}
func SRV(b ServiceBackend, zone string, state request.Request, opt Options) (records, extra []dns.RR, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	services, err := b.Services(state, false, opt)
	if err != nil {
		return nil, nil, err
	}
	dup := make(map[item]struct{})
	lookup := make(map[string]struct{})
	w := make(map[int]int)
	for _, serv := range services {
		weight := 100
		if serv.Weight != 0 {
			weight = serv.Weight
		}
		if _, ok := w[serv.Priority]; !ok {
			w[serv.Priority] = weight
			continue
		}
		w[serv.Priority] += weight
	}
	for _, serv := range services {
		if serv.Port == -1 {
			continue
		}
		w1 := 100.0 / float64(w[serv.Priority])
		if serv.Weight == 0 {
			w1 *= 100
		} else {
			w1 *= float64(serv.Weight)
		}
		weight := uint16(math.Floor(w1))
		what, ip := serv.HostType()
		switch what {
		case dns.TypeCNAME:
			srv := serv.NewSRV(state.QName(), weight)
			records = append(records, srv)
			if _, ok := lookup[srv.Target]; ok {
				break
			}
			lookup[srv.Target] = struct{}{}
			if !dns.IsSubDomain(zone, srv.Target) {
				m1, e1 := b.Lookup(state, srv.Target, dns.TypeA)
				if e1 == nil {
					extra = append(extra, m1.Answer...)
				}
				m1, e1 = b.Lookup(state, srv.Target, dns.TypeAAAA)
				if e1 == nil {
					for _, a := range m1.Answer {
						if _, ok := a.(*dns.CNAME); !ok {
							extra = append(extra, a)
						}
					}
				}
				break
			}
			state1 := state.NewWithQuestion(srv.Target, dns.TypeA)
			addr, e1 := A(b, zone, state1, nil, opt)
			if e1 == nil {
				extra = append(extra, addr...)
			}
		case dns.TypeA, dns.TypeAAAA:
			addr := serv.Host
			serv.Host = msg.Domain(serv.Key)
			srv := serv.NewSRV(state.QName(), weight)
			if ok := isDuplicate(dup, srv.Target, "", srv.Port); !ok {
				records = append(records, srv)
			}
			if ok := isDuplicate(dup, srv.Target, addr, 0); !ok {
				extra = append(extra, newAddress(serv, srv.Target, ip, what))
			}
		}
	}
	return records, extra, nil
}
func MX(b ServiceBackend, zone string, state request.Request, opt Options) (records, extra []dns.RR, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	services, err := b.Services(state, false, opt)
	if err != nil {
		return nil, nil, err
	}
	dup := make(map[item]struct{})
	lookup := make(map[string]struct{})
	for _, serv := range services {
		if !serv.Mail {
			continue
		}
		what, ip := serv.HostType()
		switch what {
		case dns.TypeCNAME:
			mx := serv.NewMX(state.QName())
			records = append(records, mx)
			if _, ok := lookup[mx.Mx]; ok {
				break
			}
			lookup[mx.Mx] = struct{}{}
			if !dns.IsSubDomain(zone, mx.Mx) {
				m1, e1 := b.Lookup(state, mx.Mx, dns.TypeA)
				if e1 == nil {
					extra = append(extra, m1.Answer...)
				}
				m1, e1 = b.Lookup(state, mx.Mx, dns.TypeAAAA)
				if e1 == nil {
					for _, a := range m1.Answer {
						if _, ok := a.(*dns.CNAME); !ok {
							extra = append(extra, a)
						}
					}
				}
				break
			}
			state1 := state.NewWithQuestion(mx.Mx, dns.TypeA)
			addr, e1 := A(b, zone, state1, nil, opt)
			if e1 == nil {
				extra = append(extra, addr...)
			}
		case dns.TypeA, dns.TypeAAAA:
			addr := serv.Host
			serv.Host = msg.Domain(serv.Key)
			mx := serv.NewMX(state.QName())
			if ok := isDuplicate(dup, mx.Mx, "", mx.Preference); !ok {
				records = append(records, mx)
			}
			if ok := isDuplicate(dup, serv.Host, addr, 0); !ok {
				extra = append(extra, newAddress(serv, serv.Host, ip, what))
			}
		}
	}
	return records, extra, nil
}
func CNAME(b ServiceBackend, zone string, state request.Request, opt Options) (records []dns.RR, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	services, err := b.Services(state, true, opt)
	if err != nil {
		return nil, err
	}
	if len(services) > 0 {
		serv := services[0]
		if ip := net.ParseIP(serv.Host); ip == nil {
			records = append(records, serv.NewCNAME(state.QName(), serv.Host))
		}
	}
	return records, nil
}
func TXT(b ServiceBackend, zone string, state request.Request, opt Options) (records []dns.RR, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	services, err := b.Services(state, false, opt)
	if err != nil {
		return nil, err
	}
	for _, serv := range services {
		if serv.Text == "" {
			continue
		}
		records = append(records, serv.NewTXT(state.QName()))
	}
	return records, nil
}
func PTR(b ServiceBackend, zone string, state request.Request, opt Options) (records []dns.RR, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	services, err := b.Reverse(state, true, opt)
	if err != nil {
		return nil, err
	}
	dup := make(map[string]struct{})
	for _, serv := range services {
		if ip := net.ParseIP(serv.Host); ip == nil {
			if _, ok := dup[serv.Host]; !ok {
				dup[serv.Host] = struct{}{}
				records = append(records, serv.NewPTR(state.QName(), serv.Host))
			}
		}
	}
	return records, nil
}
func NS(b ServiceBackend, zone string, state request.Request, opt Options) (records, extra []dns.RR, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	old := state.QName()
	state.Clear()
	state.Req.Question[0].Name = "ns.dns." + zone
	services, err := b.Services(state, false, opt)
	if err != nil {
		return nil, nil, err
	}
	state.Req.Question[0].Name = old
	for _, serv := range services {
		what, ip := serv.HostType()
		switch what {
		case dns.TypeCNAME:
			return nil, nil, fmt.Errorf("NS record must be an IP address: %s", serv.Host)
		case dns.TypeA, dns.TypeAAAA:
			serv.Host = msg.Domain(serv.Key)
			records = append(records, serv.NewNS(state.QName()))
			extra = append(extra, newAddress(serv, serv.Host, ip, what))
		}
	}
	return records, extra, nil
}
func SOA(b ServiceBackend, zone string, state request.Request, opt Options) ([]dns.RR, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	minTTL := b.MinTTL(state)
	ttl := uint32(300)
	if minTTL < ttl {
		ttl = minTTL
	}
	header := dns.RR_Header{Name: zone, Rrtype: dns.TypeSOA, Ttl: ttl, Class: dns.ClassINET}
	Mbox := hostmaster + "."
	Ns := "ns.dns."
	if zone[0] != '.' {
		Mbox += zone
		Ns += zone
	}
	soa := &dns.SOA{Hdr: header, Mbox: Mbox, Ns: Ns, Serial: b.Serial(state), Refresh: 7200, Retry: 1800, Expire: 86400, Minttl: minTTL}
	return []dns.RR{soa}, nil
}
func BackendError(b ServiceBackend, zone string, rcode int, state request.Request, err error, opt Options) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := new(dns.Msg)
	m.SetRcode(state.Req, rcode)
	m.Authoritative = true
	m.Ns, _ = SOA(b, zone, state, opt)
	state.W.WriteMsg(m)
	return dns.RcodeSuccess, err
}
func newAddress(s msg.Service, name string, ip net.IP, what uint16) dns.RR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	hdr := dns.RR_Header{Name: name, Rrtype: what, Class: dns.ClassINET, Ttl: s.TTL}
	if what == dns.TypeA {
		return &dns.A{Hdr: hdr, A: ip}
	}
	return &dns.AAAA{Hdr: hdr, AAAA: ip}
}
func checkForApex(b ServiceBackend, zone string, state request.Request, opt Options) ([]msg.Service, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if state.Name() != zone {
		return b.Services(state, false, opt)
	}
	old := state.QName()
	state.Clear()
	state.Req.Question[0].Name = dnsutil.Join("apex.dns", zone)
	services, err := b.Services(state, false, opt)
	if err == nil {
		state.Req.Question[0].Name = old
		return services, err
	}
	state.Req.Question[0].Name = old
	return b.Services(state, false, opt)
}

type item struct {
	name	string
	port	uint16
	addr	string
}

func isDuplicate(m map[item]struct{}, name, addr string, port uint16) bool {
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

const hostmaster = "hostmaster"
