package file

import (
	"context"
	"fmt"
	"io"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

var log = clog.NewWithPlugin("file")

type (
	File	struct {
		Next	plugin.Handler
		Zones	Zones
	}
	Zones	struct {
		Z	map[string]*Zone
		Names	[]string
	}
)

func (f File) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	state := request.Request{W: w, Req: r, Context: ctx}
	qname := state.Name()
	zone := plugin.Zones(f.Zones.Names).Matches(qname)
	if zone == "" {
		return plugin.NextOrFailure(f.Name(), f.Next, ctx, w, r)
	}
	z, ok := f.Zones.Z[zone]
	if !ok || z == nil {
		return dns.RcodeServerFailure, nil
	}
	if r.Opcode == dns.OpcodeNotify {
		if z.isNotify(state) {
			m := new(dns.Msg)
			m.SetReply(r)
			m.Authoritative = true
			w.WriteMsg(m)
			log.Infof("Notify from %s for %s: checking transfer", state.IP(), zone)
			ok, err := z.shouldTransfer()
			if ok {
				z.TransferIn()
			} else {
				log.Infof("Notify from %s for %s: no serial increase seen", state.IP(), zone)
			}
			if err != nil {
				log.Warningf("Notify from %s for %s: failed primary check: %s", state.IP(), zone, err)
			}
			return dns.RcodeSuccess, nil
		}
		log.Infof("Dropping notify from %s for %s", state.IP(), zone)
		return dns.RcodeSuccess, nil
	}
	if z.Expired != nil && *z.Expired {
		log.Errorf("Zone %s is expired", zone)
		return dns.RcodeServerFailure, nil
	}
	if state.QType() == dns.TypeAXFR || state.QType() == dns.TypeIXFR {
		xfr := Xfr{z}
		return xfr.ServeDNS(ctx, w, r)
	}
	answer, ns, extra, result := z.Lookup(state, qname)
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer, m.Ns, m.Extra = answer, ns, extra
	switch result {
	case Success:
	case NoData:
	case NameError:
		m.Rcode = dns.RcodeNameError
	case Delegation:
		m.Authoritative = false
	case ServerFailure:
		return dns.RcodeServerFailure, nil
	}
	w.WriteMsg(m)
	return dns.RcodeSuccess, nil
}
func (f File) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "file"
}

type serialErr struct {
	err	string
	zone	string
	origin	string
	serial	int64
}

func (s *serialErr) Error() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%s for origin %s in file %s, with serial %d", s.err, s.origin, s.zone, s.serial)
}
func Parse(f io.Reader, origin, fileName string, serial int64) (*Zone, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	zp := dns.NewZoneParser(f, dns.Fqdn(origin), fileName)
	zp.SetIncludeAllowed(true)
	z := NewZone(origin, fileName)
	seenSOA := false
	for rr, ok := zp.Next(); ok; rr, ok = zp.Next() {
		if err := zp.Err(); err != nil {
			return nil, err
		}
		if !seenSOA && serial >= 0 {
			if s, ok := rr.(*dns.SOA); ok {
				if s.Serial == uint32(serial) {
					return nil, &serialErr{err: "no change in SOA serial", origin: origin, zone: fileName, serial: serial}
				}
				seenSOA = true
			}
		}
		if err := z.Insert(rr); err != nil {
			return nil, err
		}
	}
	if !seenSOA {
		return nil, fmt.Errorf("file %q has no SOA record", fileName)
	}
	return z, nil
}
