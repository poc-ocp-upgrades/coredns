package proxy

import (
	"net"
	"github.com/miekg/dns"
)

type fakeBootWriter struct{ dns.ResponseWriter }

func (w *fakeBootWriter) LocalAddr() net.Addr {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	local := net.ParseIP("127.0.0.1")
	return &net.UDPAddr{IP: local, Port: 53}
}
func (w *fakeBootWriter) RemoteAddr() net.Addr {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	remote := net.ParseIP("8.8.8.8")
	return &net.UDPAddr{IP: remote, Port: 53}
}
