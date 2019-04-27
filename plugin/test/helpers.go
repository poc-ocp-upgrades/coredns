package test

import (
	"context"
	"sort"
	"testing"
	"github.com/miekg/dns"
)

type sect int

const (
	Answer	sect	= iota
	Ns
	Extra
)

type RRSet []dns.RR

func (p RRSet) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(p)
}
func (p RRSet) Swap(i, j int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	p[i], p[j] = p[j], p[i]
}
func (p RRSet) Less(i, j int) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return p[i].String() < p[j].String()
}

type Case struct {
	Qname	string
	Qtype	uint16
	Rcode	int
	Do	bool
	Answer	[]dns.RR
	Ns	[]dns.RR
	Extra	[]dns.RR
	Error	error
}

func (c Case) Msg() *dns.Msg {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(c.Qname), c.Qtype)
	if c.Do {
		o := new(dns.OPT)
		o.Hdr.Name = "."
		o.Hdr.Rrtype = dns.TypeOPT
		o.SetDo()
		o.SetUDPSize(4096)
		m.Extra = []dns.RR{o}
	}
	return m
}
func A(rr string) *dns.A {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r, _ := dns.NewRR(rr)
	return r.(*dns.A)
}
func AAAA(rr string) *dns.AAAA {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r, _ := dns.NewRR(rr)
	return r.(*dns.AAAA)
}
func CNAME(rr string) *dns.CNAME {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r, _ := dns.NewRR(rr)
	return r.(*dns.CNAME)
}
func DNAME(rr string) *dns.DNAME {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r, _ := dns.NewRR(rr)
	return r.(*dns.DNAME)
}
func SRV(rr string) *dns.SRV {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r, _ := dns.NewRR(rr)
	return r.(*dns.SRV)
}
func SOA(rr string) *dns.SOA {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r, _ := dns.NewRR(rr)
	return r.(*dns.SOA)
}
func NS(rr string) *dns.NS {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r, _ := dns.NewRR(rr)
	return r.(*dns.NS)
}
func PTR(rr string) *dns.PTR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r, _ := dns.NewRR(rr)
	return r.(*dns.PTR)
}
func TXT(rr string) *dns.TXT {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r, _ := dns.NewRR(rr)
	return r.(*dns.TXT)
}
func HINFO(rr string) *dns.HINFO {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r, _ := dns.NewRR(rr)
	return r.(*dns.HINFO)
}
func MX(rr string) *dns.MX {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r, _ := dns.NewRR(rr)
	return r.(*dns.MX)
}
func RRSIG(rr string) *dns.RRSIG {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r, _ := dns.NewRR(rr)
	return r.(*dns.RRSIG)
}
func NSEC(rr string) *dns.NSEC {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r, _ := dns.NewRR(rr)
	return r.(*dns.NSEC)
}
func DNSKEY(rr string) *dns.DNSKEY {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r, _ := dns.NewRR(rr)
	return r.(*dns.DNSKEY)
}
func DS(rr string) *dns.DS {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r, _ := dns.NewRR(rr)
	return r.(*dns.DS)
}
func OPT(bufsize int, do bool) *dns.OPT {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	o := new(dns.OPT)
	o.Hdr.Name = "."
	o.Hdr.Rrtype = dns.TypeOPT
	o.SetVersion(0)
	o.SetUDPSize(uint16(bufsize))
	if do {
		o.SetDo()
	}
	return o
}
func Header(t *testing.T, tc Case, resp *dns.Msg) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if resp.Rcode != tc.Rcode {
		t.Errorf("Rcode is %q, expected %q", dns.RcodeToString[resp.Rcode], dns.RcodeToString[tc.Rcode])
		return false
	}
	if len(resp.Answer) != len(tc.Answer) {
		t.Errorf("Answer for %q contained %d results, %d expected", tc.Qname, len(resp.Answer), len(tc.Answer))
		return false
	}
	if len(resp.Ns) != len(tc.Ns) {
		t.Errorf("Authority for %q contained %d results, %d expected", tc.Qname, len(resp.Ns), len(tc.Ns))
		return false
	}
	if len(resp.Extra) != len(tc.Extra) {
		t.Errorf("Additional for %q contained %d results, %d expected", tc.Qname, len(resp.Extra), len(tc.Extra))
		return false
	}
	return true
}
func Section(t *testing.T, tc Case, sec sect, rr []dns.RR) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	section := []dns.RR{}
	switch sec {
	case 0:
		section = tc.Answer
	case 1:
		section = tc.Ns
	case 2:
		section = tc.Extra
	}
	for i, a := range rr {
		if a.Header().Name != section[i].Header().Name {
			t.Errorf("RR %d should have a Header Name of %q, but has %q", i, section[i].Header().Name, a.Header().Name)
			return false
		}
		if section[i].Header().Ttl != 303 && a.Header().Ttl != section[i].Header().Ttl {
			if _, ok := section[i].(*dns.OPT); !ok {
				t.Errorf("RR %d should have a Header TTL of %d, but has %d", i, section[i].Header().Ttl, a.Header().Ttl)
				return false
			}
		}
		if a.Header().Rrtype != section[i].Header().Rrtype {
			t.Errorf("RR %d should have a header rr type of %d, but has %d", i, section[i].Header().Rrtype, a.Header().Rrtype)
			return false
		}
		switch x := a.(type) {
		case *dns.SRV:
			if x.Priority != section[i].(*dns.SRV).Priority {
				t.Errorf("RR %d should have a Priority of %d, but has %d", i, section[i].(*dns.SRV).Priority, x.Priority)
				return false
			}
			if x.Weight != section[i].(*dns.SRV).Weight {
				t.Errorf("RR %d should have a Weight of %d, but has %d", i, section[i].(*dns.SRV).Weight, x.Weight)
				return false
			}
			if x.Port != section[i].(*dns.SRV).Port {
				t.Errorf("RR %d should have a Port of %d, but has %d", i, section[i].(*dns.SRV).Port, x.Port)
				return false
			}
			if x.Target != section[i].(*dns.SRV).Target {
				t.Errorf("RR %d should have a Target of %q, but has %q", i, section[i].(*dns.SRV).Target, x.Target)
				return false
			}
		case *dns.RRSIG:
			if x.TypeCovered != section[i].(*dns.RRSIG).TypeCovered {
				t.Errorf("RR %d should have a TypeCovered of %d, but has %d", i, section[i].(*dns.RRSIG).TypeCovered, x.TypeCovered)
				return false
			}
			if x.Labels != section[i].(*dns.RRSIG).Labels {
				t.Errorf("RR %d should have a Labels of %d, but has %d", i, section[i].(*dns.RRSIG).Labels, x.Labels)
				return false
			}
			if x.SignerName != section[i].(*dns.RRSIG).SignerName {
				t.Errorf("RR %d should have a SignerName of %s, but has %s", i, section[i].(*dns.RRSIG).SignerName, x.SignerName)
				return false
			}
		case *dns.NSEC:
			if x.NextDomain != section[i].(*dns.NSEC).NextDomain {
				t.Errorf("RR %d should have a NextDomain of %s, but has %s", i, section[i].(*dns.NSEC).NextDomain, x.NextDomain)
				return false
			}
		case *dns.A:
			if x.A.String() != section[i].(*dns.A).A.String() {
				t.Errorf("RR %d should have a Address of %q, but has %q", i, section[i].(*dns.A).A.String(), x.A.String())
				return false
			}
		case *dns.AAAA:
			if x.AAAA.String() != section[i].(*dns.AAAA).AAAA.String() {
				t.Errorf("RR %d should have a Address of %q, but has %q", i, section[i].(*dns.AAAA).AAAA.String(), x.AAAA.String())
				return false
			}
		case *dns.TXT:
			for j, txt := range x.Txt {
				if txt != section[i].(*dns.TXT).Txt[j] {
					t.Errorf("RR %d should have a Txt of %q, but has %q", i, section[i].(*dns.TXT).Txt[j], txt)
					return false
				}
			}
		case *dns.HINFO:
			if x.Cpu != section[i].(*dns.HINFO).Cpu {
				t.Errorf("RR %d should have a Cpu of %s, but has %s", i, section[i].(*dns.HINFO).Cpu, x.Cpu)
			}
			if x.Os != section[i].(*dns.HINFO).Os {
				t.Errorf("RR %d should have a Os of %s, but has %s", i, section[i].(*dns.HINFO).Os, x.Os)
			}
		case *dns.SOA:
			tt := section[i].(*dns.SOA)
			if x.Ns != tt.Ns {
				t.Errorf("SOA nameserver should be %q, but is %q", tt.Ns, x.Ns)
				return false
			}
		case *dns.PTR:
			tt := section[i].(*dns.PTR)
			if x.Ptr != tt.Ptr {
				t.Errorf("PTR ptr should be %q, but is %q", tt.Ptr, x.Ptr)
				return false
			}
		case *dns.CNAME:
			tt := section[i].(*dns.CNAME)
			if x.Target != tt.Target {
				t.Errorf("CNAME target should be %q, but is %q", tt.Target, x.Target)
				return false
			}
		case *dns.MX:
			tt := section[i].(*dns.MX)
			if x.Mx != tt.Mx {
				t.Errorf("MX Mx should be %q, but is %q", tt.Mx, x.Mx)
				return false
			}
			if x.Preference != tt.Preference {
				t.Errorf("MX Preference should be %q, but is %q", tt.Preference, x.Preference)
				return false
			}
		case *dns.NS:
			tt := section[i].(*dns.NS)
			if x.Ns != tt.Ns {
				t.Errorf("NS nameserver should be %q, but is %q", tt.Ns, x.Ns)
				return false
			}
		case *dns.OPT:
			tt := section[i].(*dns.OPT)
			if x.UDPSize() != tt.UDPSize() {
				t.Errorf("OPT UDPSize should be %d, but is %d", tt.UDPSize(), x.UDPSize())
				return false
			}
			if x.Do() != tt.Do() {
				t.Errorf("OPT DO should be %t, but is %t", tt.Do(), x.Do())
				return false
			}
		}
	}
	return true
}
func CNAMEOrder(t *testing.T, res *dns.Msg) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for i, c := range res.Answer {
		if c.Header().Rrtype != dns.TypeCNAME {
			continue
		}
		for _, a := range res.Answer[:i] {
			if a.Header().Name != c.(*dns.CNAME).Target {
				continue
			}
			t.Errorf("CNAME found after target record\n")
			t.Logf("%v\n", res)
		}
	}
}
func SortAndCheck(t *testing.T, resp *dns.Msg, tc Case) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	sort.Sort(RRSet(resp.Answer))
	sort.Sort(RRSet(resp.Ns))
	sort.Sort(RRSet(resp.Extra))
	if !Header(t, tc, resp) {
		t.Logf("%v\n", resp)
		return
	}
	if !Section(t, tc, Answer, resp.Answer) {
		t.Logf("%v\n", resp)
		return
	}
	if !Section(t, tc, Ns, resp.Ns) {
		t.Logf("%v\n", resp)
		return
	}
	if !Section(t, tc, Extra, resp.Extra) {
		t.Logf("%v\n", resp)
		return
	}
	return
}
func ErrorHandler() Handler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return HandlerFunc(func(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
		m := new(dns.Msg)
		m.SetRcode(r, dns.RcodeServerFailure)
		w.WriteMsg(m)
		return dns.RcodeServerFailure, nil
	})
}
func NextHandler(rcode int, err error) Handler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return HandlerFunc(func(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
		return rcode, err
	})
}

type (
	HandlerFunc	func(context.Context, dns.ResponseWriter, *dns.Msg) (int, error)
	Handler		interface {
		ServeDNS(context.Context, dns.ResponseWriter, *dns.Msg) (int, error)
		Name() string
	}
)

func (f HandlerFunc) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return f(ctx, w, r)
}
func (f HandlerFunc) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "handlerfunc"
}
