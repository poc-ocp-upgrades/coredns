package dnssec

import (
	"time"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/cache"
	"github.com/coredns/coredns/plugin/pkg/response"
	"github.com/coredns/coredns/plugin/pkg/singleflight"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Dnssec struct {
	Next		plugin.Handler
	zones		[]string
	keys		[]*DNSKEY
	splitkeys	bool
	inflight	*singleflight.Group
	cache		*cache.Cache
}

func New(zones []string, keys []*DNSKEY, splitkeys bool, next plugin.Handler, c *cache.Cache) Dnssec {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return Dnssec{Next: next, zones: zones, keys: keys, splitkeys: splitkeys, cache: c, inflight: new(singleflight.Group)}
}
func (d Dnssec) Sign(state request.Request, now time.Time, server string) *dns.Msg {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req := state.Req
	incep, expir := incepExpir(now)
	mt, _ := response.Typify(req, time.Now().UTC())
	if mt == response.Delegation {
		return req
	}
	if mt == response.NameError || mt == response.NoData {
		if req.Ns[0].Header().Rrtype != dns.TypeSOA || len(req.Ns) > 1 {
			return req
		}
		ttl := req.Ns[0].Header().Ttl
		if sigs, err := d.sign(req.Ns, state.Zone, ttl, incep, expir, server); err == nil {
			req.Ns = append(req.Ns, sigs...)
		}
		if sigs, err := d.nsec(state, mt, ttl, incep, expir, server); err == nil {
			req.Ns = append(req.Ns, sigs...)
		}
		if len(req.Ns) > 1 {
			req.Rcode = dns.RcodeSuccess
		}
		return req
	}
	for _, r := range rrSets(req.Answer) {
		ttl := r[0].Header().Ttl
		if sigs, err := d.sign(r, state.Zone, ttl, incep, expir, server); err == nil {
			req.Answer = append(req.Answer, sigs...)
		}
	}
	for _, r := range rrSets(req.Ns) {
		ttl := r[0].Header().Ttl
		if sigs, err := d.sign(r, state.Zone, ttl, incep, expir, server); err == nil {
			req.Ns = append(req.Ns, sigs...)
		}
	}
	for _, r := range rrSets(req.Extra) {
		ttl := r[0].Header().Ttl
		if sigs, err := d.sign(r, state.Zone, ttl, incep, expir, server); err == nil {
			req.Extra = append(req.Extra, sigs...)
		}
	}
	return req
}
func (d Dnssec) sign(rrs []dns.RR, signerName string, ttl, incep, expir uint32, server string) ([]dns.RR, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	k := hash(rrs)
	sgs, ok := d.get(k, server)
	if ok {
		return sgs, nil
	}
	sigs, err := d.inflight.Do(k, func() (interface{}, error) {
		var sigs []dns.RR
		for _, k := range d.keys {
			if d.splitkeys {
				if len(rrs) > 0 && rrs[0].Header().Rrtype == dns.TypeDNSKEY {
					if !k.isKSK() {
						continue
					}
				} else {
					if !k.isZSK() {
						continue
					}
				}
			}
			sig := k.newRRSIG(signerName, ttl, incep, expir)
			if e := sig.Sign(k.s, rrs); e != nil {
				return sigs, e
			}
			sigs = append(sigs, sig)
		}
		d.set(k, sigs)
		return sigs, nil
	})
	return sigs.([]dns.RR), err
}
func (d Dnssec) set(key uint64, sigs []dns.RR) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	d.cache.Add(key, sigs)
}
func (d Dnssec) get(key uint64, server string) ([]dns.RR, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s, ok := d.cache.Get(key); ok {
		is75 := time.Now().UTC().Add(sixDays)
		for _, rr := range s.([]dns.RR) {
			if !rr.(*dns.RRSIG).ValidityPeriod(is75) {
				cacheMisses.WithLabelValues(server).Inc()
				return nil, false
			}
		}
		cacheHits.WithLabelValues(server).Inc()
		return s.([]dns.RR), true
	}
	cacheMisses.WithLabelValues(server).Inc()
	return nil, false
}
func incepExpir(now time.Time) (uint32, uint32) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	incep := uint32(now.Add(-3 * time.Hour).Unix())
	expir := uint32(now.Add(eightDays).Unix())
	return incep, expir
}

const (
	eightDays	= 8 * 24 * time.Hour
	sixDays		= 6 * 24 * time.Hour
	defaultCap	= 10000
)
