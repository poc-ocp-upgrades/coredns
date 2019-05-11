package dnsserver

import (
	"net"
	"github.com/coredns/coredns/plugin/pkg/nonwriter"
)

type DoHWriter struct {
	nonwriter.Writer
	raddr	net.Addr
	laddr	net.Addr
}

func (d *DoHWriter) RemoteAddr() net.Addr {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return d.raddr
}
func (d *DoHWriter) LocalAddr() net.Addr {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return d.laddr
}
