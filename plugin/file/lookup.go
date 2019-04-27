package file

import (
	"github.com/coredns/coredns/plugin/file/tree"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Result int

const (
	Success	Result	= iota
	NameError
	Delegation
	NoData
	ServerFailure
)

func (z *Zone) Lookup(state request.Request, qname string) ([]dns.RR, []dns.RR, []dns.RR, Result) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	qtype := state.QType()
	do := state.Do()
	if 0 < z.ReloadInterval {
		z.reloadMu.RLock()
	}
	defer func() {
		if 0 < z.ReloadInterval {
			z.reloadMu.RUnlock()
		}
	}()
	soa := z.Apex.SOA
	if soa == nil {
		return nil, nil, nil, ServerFailure
	}
	if qtype == dns.TypeSOA {
		return z.soa(do), z.ns(do), nil, Success
	}
	if qtype == dns.TypeNS && qname == z.origin {
		nsrrs := z.ns(do)
		glue := z.Glue(nsrrs, do)
		return nsrrs, nil, glue, Success
	}
	var (
		found, shot	bool
		parts		string
		i		int
		elem, wildElem	*tree.Elem
	)
	for {
		parts, shot = z.nameFromRight(qname, i)
		if shot {
			break
		}
		elem, found = z.Tree.Search(parts)
		if !found {
			wildcard := replaceWithAsteriskLabel(parts)
			if wild, found := z.Tree.Search(wildcard); found {
				wildElem = wild
			}
			i++
			continue
		}
		if dnamerrs := elem.Types(dns.TypeDNAME); dnamerrs != nil {
			dname := dnamerrs[0]
			if cname := synthesizeCNAME(state.Name(), dname.(*dns.DNAME)); cname != nil {
				answer, ns, extra, rcode := z.additionalProcessing(state, elem, []dns.RR{cname})
				if do {
					sigs := elem.Types(dns.TypeRRSIG)
					sigs = signatureForSubType(sigs, dns.TypeDNAME)
					dnamerrs = append(dnamerrs, sigs...)
				}
				answer = append(dnamerrs, answer...)
				return answer, ns, extra, rcode
			}
		}
		if nsrrs := elem.Types(dns.TypeNS); nsrrs != nil {
			if qtype == dns.TypeDS && elem.Name() == qname {
				i++
				continue
			}
			glue := z.Glue(nsrrs, do)
			if do {
				dss := z.typeFromElem(elem, dns.TypeDS, do)
				nsrrs = append(nsrrs, dss...)
			}
			return nil, nsrrs, glue, Delegation
		}
		i++
	}
	if found && !shot {
		return nil, nil, nil, ServerFailure
	}
	if found && shot {
		if rrs := elem.Types(dns.TypeCNAME); len(rrs) > 0 && qtype != dns.TypeCNAME {
			return z.additionalProcessing(state, elem, rrs)
		}
		rrs := elem.Types(qtype, qname)
		if len(rrs) == 0 {
			ret := z.soa(do)
			if do {
				nsec := z.typeFromElem(elem, dns.TypeNSEC, do)
				ret = append(ret, nsec...)
			}
			return nil, ret, nil, NoData
		}
		additional := additionalProcessing(z, rrs, do)
		if do {
			sigs := elem.Types(dns.TypeRRSIG)
			sigs = signatureForSubType(sigs, qtype)
			rrs = append(rrs, sigs...)
		}
		return rrs, z.ns(do), additional, Success
	}
	if wildElem != nil {
		auth := z.ns(do)
		if rrs := wildElem.Types(dns.TypeCNAME, qname); len(rrs) > 0 {
			return z.additionalProcessing(state, wildElem, rrs)
		}
		rrs := wildElem.Types(qtype, qname)
		if len(rrs) == 0 {
			ret := z.soa(do)
			if do {
				nsec := z.typeFromElem(wildElem, dns.TypeNSEC, do)
				ret = append(ret, nsec...)
			}
			return nil, ret, nil, Success
		}
		if do {
			if deny, found := z.Tree.Prev(qname); found {
				nsec := z.typeFromElem(deny, dns.TypeNSEC, do)
				auth = append(auth, nsec...)
			}
			sigs := wildElem.Types(dns.TypeRRSIG, qname)
			sigs = signatureForSubType(sigs, qtype)
			rrs = append(rrs, sigs...)
		}
		return rrs, auth, nil, Success
	}
	rcode := NameError
	if x, found := z.Tree.Next(qname); found {
		if dns.IsSubDomain(qname, x.Name()) {
			rcode = Success
		}
	}
	ret := z.soa(do)
	if do {
		deny, found := z.Tree.Prev(qname)
		if !found {
			goto Out
		}
		nsec := z.typeFromElem(deny, dns.TypeNSEC, do)
		ret = append(ret, nsec...)
		if rcode != NameError {
			goto Out
		}
		ce, found := z.ClosestEncloser(qname)
		if found {
			wildcard := "*." + ce.Name()
			if ss, found := z.Tree.Prev(wildcard); found {
				if ss.Name() != deny.Name() {
					nsec := z.typeFromElem(ss, dns.TypeNSEC, do)
					ret = append(ret, nsec...)
				}
			}
		}
	}
Out:
	return nil, ret, nil, rcode
}
func (z *Zone) typeFromElem(elem *tree.Elem, tp uint16, do bool) []dns.RR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	rrs := elem.Types(tp)
	if do {
		sigs := elem.Types(dns.TypeRRSIG)
		sigs = signatureForSubType(sigs, tp)
		if len(sigs) > 0 {
			rrs = append(rrs, sigs...)
		}
	}
	return rrs
}
func (z *Zone) soa(do bool) []dns.RR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if do {
		ret := append([]dns.RR{z.Apex.SOA}, z.Apex.SIGSOA...)
		return ret
	}
	return []dns.RR{z.Apex.SOA}
}
func (z *Zone) ns(do bool) []dns.RR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if do {
		ret := append(z.Apex.NS, z.Apex.SIGNS...)
		return ret
	}
	return z.Apex.NS
}
func (z *Zone) additionalProcessing(state request.Request, elem *tree.Elem, rrs []dns.RR) ([]dns.RR, []dns.RR, []dns.RR, Result) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	qtype := state.QType()
	do := state.Do()
	if do {
		sigs := elem.Types(dns.TypeRRSIG)
		sigs = signatureForSubType(sigs, dns.TypeCNAME)
		if len(sigs) > 0 {
			rrs = append(rrs, sigs...)
		}
	}
	targetName := rrs[0].(*dns.CNAME).Target
	elem, _ = z.Tree.Search(targetName)
	if elem == nil {
		rrs = append(rrs, z.externalLookup(state, targetName, qtype)...)
		return rrs, z.ns(do), nil, Success
	}
	i := 0
Redo:
	cname := elem.Types(dns.TypeCNAME)
	if len(cname) > 0 {
		rrs = append(rrs, cname...)
		if do {
			sigs := elem.Types(dns.TypeRRSIG)
			sigs = signatureForSubType(sigs, dns.TypeCNAME)
			if len(sigs) > 0 {
				rrs = append(rrs, sigs...)
			}
		}
		targetName := cname[0].(*dns.CNAME).Target
		elem, _ = z.Tree.Search(targetName)
		if elem == nil {
			rrs = append(rrs, z.externalLookup(state, targetName, qtype)...)
			return rrs, z.ns(do), nil, Success
		}
		i++
		if i > maxChain {
			return rrs, z.ns(do), nil, Success
		}
		goto Redo
	}
	targets := cnameForType(elem.All(), qtype)
	if len(targets) > 0 {
		rrs = append(rrs, targets...)
		if do {
			sigs := elem.Types(dns.TypeRRSIG)
			sigs = signatureForSubType(sigs, qtype)
			if len(sigs) > 0 {
				rrs = append(rrs, sigs...)
			}
		}
	}
	return rrs, z.ns(do), nil, Success
}
func cnameForType(targets []dns.RR, origQtype uint16) []dns.RR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ret := []dns.RR{}
	for _, target := range targets {
		if target.Header().Rrtype == origQtype {
			ret = append(ret, target)
		}
	}
	return ret
}
func (z *Zone) externalLookup(state request.Request, target string, qtype uint16) []dns.RR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m, e := z.Upstream.Lookup(state, target, qtype)
	if e != nil {
		return nil
	}
	if m == nil {
		return nil
	}
	return m.Answer
}
func signatureForSubType(rrs []dns.RR, subtype uint16) []dns.RR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	sigs := []dns.RR{}
	for _, sig := range rrs {
		if s, ok := sig.(*dns.RRSIG); ok {
			if s.TypeCovered == subtype {
				sigs = append(sigs, s)
			}
		}
	}
	return sigs
}
func (z *Zone) Glue(nsrrs []dns.RR, do bool) []dns.RR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	glue := []dns.RR{}
	for _, rr := range nsrrs {
		if ns, ok := rr.(*dns.NS); ok && dns.IsSubDomain(ns.Header().Name, ns.Ns) {
			glue = append(glue, z.searchGlue(ns.Ns, do)...)
		}
	}
	return glue
}
func (z *Zone) searchGlue(name string, do bool) []dns.RR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	glue := []dns.RR{}
	if elem, found := z.Tree.Search(name); found {
		glue = append(glue, elem.Types(dns.TypeA)...)
		if do {
			sigs := elem.Types(dns.TypeRRSIG)
			sigs = signatureForSubType(sigs, dns.TypeA)
			glue = append(glue, sigs...)
		}
	}
	if elem, found := z.Tree.Search(name); found {
		glue = append(glue, elem.Types(dns.TypeAAAA)...)
		if do {
			sigs := elem.Types(dns.TypeRRSIG)
			sigs = signatureForSubType(sigs, dns.TypeAAAA)
			glue = append(glue, sigs...)
		}
	}
	return glue
}
func additionalProcessing(z *Zone, answer []dns.RR, do bool) (extra []dns.RR) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, rr := range answer {
		name := ""
		switch x := rr.(type) {
		case *dns.SRV:
			name = x.Target
		case *dns.MX:
			name = x.Mx
		}
		if !dns.IsSubDomain(z.origin, name) {
			continue
		}
		elem, _ := z.Tree.Search(name)
		if elem == nil {
			continue
		}
		sigs := elem.Types(dns.TypeRRSIG)
		for _, addr := range []uint16{dns.TypeA, dns.TypeAAAA} {
			if a := elem.Types(addr); a != nil {
				extra = append(extra, a...)
				if do {
					sig := signatureForSubType(sigs, addr)
					extra = append(extra, sig...)
				}
			}
		}
	}
	return extra
}

const maxChain = 8
