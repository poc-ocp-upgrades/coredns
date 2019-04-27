package erratic

import (
	"context"
	"sync/atomic"
	"time"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Erratic struct {
	drop		uint64
	delay		uint64
	duration	time.Duration
	truncate	uint64
	large		bool
	q		uint64
}

func (e *Erratic) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	state := request.Request{W: w, Req: r}
	drop := false
	delay := false
	trunc := false
	queryNr := atomic.LoadUint64(&e.q)
	atomic.AddUint64(&e.q, 1)
	if e.drop > 0 && queryNr%e.drop == 0 {
		drop = true
	}
	if e.delay > 0 && queryNr%e.delay == 0 {
		delay = true
	}
	if e.truncate > 0 && queryNr&e.truncate == 0 {
		trunc = true
	}
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	if trunc {
		m.Truncated = true
	}
	switch state.QType() {
	case dns.TypeA:
		rr := *(rrA.(*dns.A))
		rr.Header().Name = state.QName()
		m.Answer = append(m.Answer, &rr)
		if e.large {
			for i := 0; i < 29; i++ {
				m.Answer = append(m.Answer, &rr)
			}
		}
	case dns.TypeAAAA:
		rr := *(rrAAAA.(*dns.AAAA))
		rr.Header().Name = state.QName()
		m.Answer = append(m.Answer, &rr)
	case dns.TypeAXFR:
		if drop {
			return 0, nil
		}
		if delay {
			time.Sleep(e.duration)
		}
		xfr(state, trunc)
		return 0, nil
	default:
		if drop {
			return 0, nil
		}
		if delay {
			time.Sleep(e.duration)
		}
		return dns.RcodeServerFailure, nil
	}
	if drop {
		return 0, nil
	}
	if delay {
		time.Sleep(e.duration)
	}
	w.WriteMsg(m)
	return 0, nil
}
func (e *Erratic) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "erratic"
}

var (
	rrA, _		= dns.NewRR(". IN 0 A 192.0.2.53")
	rrAAAA, _	= dns.NewRR(". IN 0 AAAA 2001:DB8::53")
)
