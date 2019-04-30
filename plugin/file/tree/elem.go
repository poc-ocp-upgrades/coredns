package tree

import "github.com/miekg/dns"

type Elem struct {
	m	map[uint16][]dns.RR
	name	string
}

func newElem(rr dns.RR) *Elem {
	_logClusterCodePath()
	defer _logClusterCodePath()
	e := Elem{m: make(map[uint16][]dns.RR)}
	e.m[rr.Header().Rrtype] = []dns.RR{rr}
	return &e
}
func (e *Elem) Types(qtype uint16, qname ...string) []dns.RR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rrs := e.m[qtype]
	if rrs != nil && len(qname) > 0 {
		copied := make([]dns.RR, len(rrs))
		for i := range rrs {
			copied[i] = dns.Copy(rrs[i])
			copied[i].Header().Name = qname[0]
		}
		return copied
	}
	return rrs
}
func (e *Elem) All() []dns.RR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	list := []dns.RR{}
	for _, rrs := range e.m {
		list = append(list, rrs...)
	}
	return list
}
func (e *Elem) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if e.name != "" {
		return e.name
	}
	for _, rrs := range e.m {
		e.name = rrs[0].Header().Name
		return e.name
	}
	return ""
}
func (e *Elem) Empty() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(e.m) == 0
}
func (e *Elem) Insert(rr dns.RR) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	t := rr.Header().Rrtype
	if e.m == nil {
		e.m = make(map[uint16][]dns.RR)
		e.m[t] = []dns.RR{rr}
		return
	}
	rrs, ok := e.m[t]
	if !ok {
		e.m[t] = []dns.RR{rr}
		return
	}
	for _, er := range rrs {
		if equalRdata(er, rr) {
			return
		}
	}
	rrs = append(rrs, rr)
	e.m[t] = rrs
}
func (e *Elem) Delete(rr dns.RR) (empty bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if e.m == nil {
		return true
	}
	t := rr.Header().Rrtype
	rrs, ok := e.m[t]
	if !ok {
		return
	}
	for i, er := range rrs {
		if equalRdata(er, rr) {
			rrs = removeFromSlice(rrs, i)
			e.m[t] = rrs
			empty = len(rrs) == 0
			if empty {
				delete(e.m, t)
			}
			return
		}
	}
	return
}
func Less(a *Elem, name string) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return less(name, a.Name())
}
func equalRdata(a, b dns.RR) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch x := a.(type) {
	case *dns.A:
		return x.A.Equal(b.(*dns.A).A)
	case *dns.AAAA:
		return x.AAAA.Equal(b.(*dns.AAAA).AAAA)
	case *dns.MX:
		if x.Mx == b.(*dns.MX).Mx && x.Preference == b.(*dns.MX).Preference {
			return true
		}
	}
	return false
}
func removeFromSlice(rrs []dns.RR, i int) []dns.RR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if i >= len(rrs) {
		return rrs
	}
	rrs = append(rrs[:i], rrs[i+1:]...)
	return rrs
}
