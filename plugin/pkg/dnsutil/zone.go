package dnsutil

import (
	"errors"
	"github.com/miekg/dns"
)

func TrimZone(q string, z string) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	zl := dns.CountLabel(z)
	i, ok := dns.PrevLabel(q, zl)
	if ok || i-1 < 0 {
		return "", errors.New("trimzone: overshot qname: " + q + "for zone " + z)
	}
	return q[:i-1], nil
}
