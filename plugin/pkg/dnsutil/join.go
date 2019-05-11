package dnsutil

import (
	"strings"
	"github.com/miekg/dns"
)

func Join(labels ...string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ll := len(labels)
	if labels[ll-1] == "." {
		return strings.Join(labels[:ll-1], ".") + "."
	}
	return dns.Fqdn(strings.Join(labels, "."))
}
