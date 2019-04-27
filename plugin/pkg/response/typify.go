package response

import (
	"fmt"
	"time"
	"github.com/miekg/dns"
)

type Type int

const (
	NoError	Type	= iota
	NameError
	NoData
	Delegation
	Meta
	Update
	OtherError
)

var toString = map[Type]string{NoError: "NOERROR", NameError: "NXDOMAIN", NoData: "NODATA", Delegation: "DELEGATION", Meta: "META", Update: "UPDATE", OtherError: "OTHERERROR"}

func (t Type) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return toString[t]
}
func TypeFromString(s string) (Type, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for t, str := range toString {
		if s == str {
			return t, nil
		}
	}
	return NoError, fmt.Errorf("invalid Type: %s", s)
}
func Typify(m *dns.Msg, t time.Time) (Type, *dns.OPT) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m == nil {
		return OtherError, nil
	}
	opt := m.IsEdns0()
	do := false
	if opt != nil {
		do = opt.Do()
	}
	if m.Opcode == dns.OpcodeUpdate {
		return Update, opt
	}
	if m.Opcode == dns.OpcodeNotify {
		return Meta, opt
	}
	if len(m.Question) > 0 {
		if m.Question[0].Qtype == dns.TypeAXFR || m.Question[0].Qtype == dns.TypeIXFR {
			return Meta, opt
		}
	}
	if do {
		if expired := typifyExpired(m, t); expired {
			return OtherError, opt
		}
	}
	if len(m.Answer) > 0 && m.Rcode == dns.RcodeSuccess {
		return NoError, opt
	}
	soa := false
	ns := 0
	for _, r := range m.Ns {
		if r.Header().Rrtype == dns.TypeSOA {
			soa = true
			continue
		}
		if r.Header().Rrtype == dns.TypeNS {
			ns++
		}
	}
	if soa && m.Rcode == dns.RcodeSuccess {
		return NoData, opt
	}
	if soa && m.Rcode == dns.RcodeNameError {
		return NameError, opt
	}
	if ns > 0 && m.Rcode == dns.RcodeSuccess {
		return Delegation, opt
	}
	if m.Rcode == dns.RcodeSuccess {
		return NoError, opt
	}
	return OtherError, opt
}
func typifyExpired(m *dns.Msg, t time.Time) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if expired := typifyExpiredRRSIG(m.Answer, t); expired {
		return true
	}
	if expired := typifyExpiredRRSIG(m.Ns, t); expired {
		return true
	}
	if expired := typifyExpiredRRSIG(m.Extra, t); expired {
		return true
	}
	return false
}
func typifyExpiredRRSIG(rrs []dns.RR, t time.Time) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, r := range rrs {
		if r.Header().Rrtype != dns.TypeRRSIG {
			continue
		}
		ok := r.(*dns.RRSIG).ValidityPeriod(t)
		if !ok {
			return true
		}
	}
	return false
}
