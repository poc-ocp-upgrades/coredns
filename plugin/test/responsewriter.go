package test

import (
	"net"
	"github.com/miekg/dns"
)

type ResponseWriter struct{ TCP bool }

func (t *ResponseWriter) LocalAddr() net.Addr {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ip := net.ParseIP("127.0.0.1")
	port := 53
	if t.TCP {
		return &net.TCPAddr{IP: ip, Port: port, Zone: ""}
	}
	return &net.UDPAddr{IP: ip, Port: port, Zone: ""}
}
func (t *ResponseWriter) RemoteAddr() net.Addr {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ip := net.ParseIP("10.240.0.1")
	port := 40212
	if t.TCP {
		return &net.TCPAddr{IP: ip, Port: port, Zone: ""}
	}
	return &net.UDPAddr{IP: ip, Port: port, Zone: ""}
}
func (t *ResponseWriter) WriteMsg(m *dns.Msg) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (t *ResponseWriter) Write(buf []byte) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(buf), nil
}
func (t *ResponseWriter) Close() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (t *ResponseWriter) TsigStatus() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (t *ResponseWriter) TsigTimersOnly(bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return
}
func (t *ResponseWriter) Hijack() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return
}

type ResponseWriter6 struct{ ResponseWriter }

func (t *ResponseWriter6) LocalAddr() net.Addr {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if t.TCP {
		return &net.TCPAddr{IP: net.ParseIP("::1"), Port: 53, Zone: ""}
	}
	return &net.UDPAddr{IP: net.ParseIP("::1"), Port: 53, Zone: ""}
}
func (t *ResponseWriter6) RemoteAddr() net.Addr {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if t.TCP {
		return &net.TCPAddr{IP: net.ParseIP("fe80::42:ff:feca:4c65"), Port: 40212, Zone: ""}
	}
	return &net.UDPAddr{IP: net.ParseIP("fe80::42:ff:feca:4c65"), Port: 40212, Zone: ""}
}
