package file

import (
	"fmt"
	"net"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"github.com/coredns/coredns/plugin/file/tree"
	"github.com/coredns/coredns/plugin/pkg/upstream"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Zone struct {
	origin	string
	origLen	int
	file	string
	*tree.Tree
	Apex			Apex
	TransferTo		[]string
	StartupOnce		sync.Once
	TransferFrom	[]string
	Expired			*bool
	ReloadInterval	time.Duration
	LastReloaded	time.Time
	reloadMu		sync.RWMutex
	reloadShutdown	chan bool
	Upstream		upstream.Upstream
}
type Apex struct {
	SOA		*dns.SOA
	NS		[]dns.RR
	SIGSOA	[]dns.RR
	SIGNS	[]dns.RR
}

func NewZone(name, file string) *Zone {
	_logClusterCodePath()
	defer _logClusterCodePath()
	z := &Zone{origin: dns.Fqdn(name), origLen: dns.CountLabel(dns.Fqdn(name)), file: filepath.Clean(file), Tree: &tree.Tree{}, Expired: new(bool), reloadShutdown: make(chan bool), LastReloaded: time.Now()}
	*z.Expired = false
	return z
}
func (z *Zone) Copy() *Zone {
	_logClusterCodePath()
	defer _logClusterCodePath()
	z1 := NewZone(z.origin, z.file)
	z1.TransferTo = z.TransferTo
	z1.TransferFrom = z.TransferFrom
	z1.Expired = z.Expired
	z1.Apex = z.Apex
	return z1
}
func (z *Zone) CopyWithoutApex() *Zone {
	_logClusterCodePath()
	defer _logClusterCodePath()
	z1 := NewZone(z.origin, z.file)
	z1.TransferTo = z.TransferTo
	z1.TransferFrom = z.TransferFrom
	z1.Expired = z.Expired
	return z1
}
func (z *Zone) Insert(r dns.RR) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.Header().Name = strings.ToLower(r.Header().Name)
	switch h := r.Header().Rrtype; h {
	case dns.TypeNS:
		r.(*dns.NS).Ns = strings.ToLower(r.(*dns.NS).Ns)
		if r.Header().Name == z.origin {
			z.Apex.NS = append(z.Apex.NS, r)
			return nil
		}
	case dns.TypeSOA:
		r.(*dns.SOA).Ns = strings.ToLower(r.(*dns.SOA).Ns)
		r.(*dns.SOA).Mbox = strings.ToLower(r.(*dns.SOA).Mbox)
		z.Apex.SOA = r.(*dns.SOA)
		return nil
	case dns.TypeNSEC3, dns.TypeNSEC3PARAM:
		return fmt.Errorf("NSEC3 zone is not supported, dropping RR: %s for zone: %s", r.Header().Name, z.origin)
	case dns.TypeRRSIG:
		x := r.(*dns.RRSIG)
		switch x.TypeCovered {
		case dns.TypeSOA:
			z.Apex.SIGSOA = append(z.Apex.SIGSOA, x)
			return nil
		case dns.TypeNS:
			if r.Header().Name == z.origin {
				z.Apex.SIGNS = append(z.Apex.SIGNS, x)
				return nil
			}
		}
	case dns.TypeCNAME:
		r.(*dns.CNAME).Target = strings.ToLower(r.(*dns.CNAME).Target)
	case dns.TypeMX:
		r.(*dns.MX).Mx = strings.ToLower(r.(*dns.MX).Mx)
	case dns.TypeSRV:
		r.(*dns.SRV).Target = strings.ToLower(r.(*dns.SRV).Target)
	}
	z.Tree.Insert(r)
	return nil
}
func (z *Zone) Delete(r dns.RR) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	z.Tree.Delete(r)
}
func (z *Zone) File() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	z.reloadMu.Lock()
	defer z.reloadMu.Unlock()
	return z.file
}
func (z *Zone) SetFile(path string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	z.reloadMu.Lock()
	z.file = path
	z.reloadMu.Unlock()
}
func (z *Zone) TransferAllowed(state request.Request) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, t := range z.TransferTo {
		if t == "*" {
			return true
		}
		remote := state.IP()
		to, _, err := net.SplitHostPort(t)
		if err != nil {
			continue
		}
		if to == remote {
			return true
		}
	}
	return false
}
func (z *Zone) All() []dns.RR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if z.ReloadInterval > 0 {
		z.reloadMu.RLock()
		defer z.reloadMu.RUnlock()
	}
	records := []dns.RR{}
	allNodes := z.Tree.All()
	for _, a := range allNodes {
		records = append(records, a.All()...)
	}
	if len(z.Apex.SIGNS) > 0 {
		records = append(z.Apex.SIGNS, records...)
	}
	records = append(z.Apex.NS, records...)
	if len(z.Apex.SIGSOA) > 0 {
		records = append(z.Apex.SIGSOA, records...)
	}
	return append([]dns.RR{z.Apex.SOA}, records...)
}
func (z *Zone) Print() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	z.Tree.Print()
}
func (z *Zone) nameFromRight(qname string, i int) (string, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if i <= 0 {
		return z.origin, false
	}
	for j := 1; j <= z.origLen; j++ {
		if _, shot := dns.PrevLabel(qname, j); shot {
			return qname, shot
		}
	}
	k := 0
	var shot bool
	for j := 1; j <= i; j++ {
		k, shot = dns.PrevLabel(qname, j+z.origLen)
		if shot {
			return qname, shot
		}
	}
	return qname[k:], false
}
