package parse

import (
	"strings"
	"github.com/coredns/coredns/plugin/pkg/transport"
)

func Transport(s string) (trans string, addr string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch {
	case strings.HasPrefix(s, transport.TLS+"://"):
		s = s[len(transport.TLS+"://"):]
		return transport.TLS, s
	case strings.HasPrefix(s, transport.DNS+"://"):
		s = s[len(transport.DNS+"://"):]
		return transport.DNS, s
	case strings.HasPrefix(s, transport.GRPC+"://"):
		s = s[len(transport.GRPC+"://"):]
		return transport.GRPC, s
	case strings.HasPrefix(s, transport.HTTPS+"://"):
		s = s[len(transport.HTTPS+"://"):]
		return transport.HTTPS, s
	}
	return transport.DNS, s
}
