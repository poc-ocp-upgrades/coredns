package dnsserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"github.com/coredns/coredns/plugin/pkg/transport"
	"github.com/miekg/dns"
)

type ServerTLS struct {
	*Server
	tlsConfig	*tls.Config
}

func NewServerTLS(addr string, group []*Config) (*ServerTLS, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s, err := NewServer(addr, group)
	if err != nil {
		return nil, err
	}
	var tlsConfig *tls.Config
	for _, conf := range s.zones {
		tlsConfig = conf.TLSConfig
	}
	return &ServerTLS{Server: s, tlsConfig: tlsConfig}, nil
}
func (s *ServerTLS) Serve(l net.Listener) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.m.Lock()
	if s.tlsConfig != nil {
		l = tls.NewListener(l, s.tlsConfig)
	}
	s.server[tcp] = &dns.Server{Listener: l, Net: "tcp-tls", Handler: dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
		ctx := context.Background()
		s.ServeDNS(ctx, w, r)
	})}
	s.m.Unlock()
	return s.server[tcp].ActivateAndServe()
}
func (s *ServerTLS) ServePacket(p net.PacketConn) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (s *ServerTLS) Listen() (net.Listener, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l, err := net.Listen("tcp", s.Addr[len(transport.TLS+"://"):])
	if err != nil {
		return nil, err
	}
	return l, nil
}
func (s *ServerTLS) ListenPacket() (net.PacketConn, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, nil
}
func (s *ServerTLS) OnStartupComplete() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if Quiet {
		return
	}
	out := startUpZones(transport.TLS+"://", s.Addr, s.zones)
	if out != "" {
		fmt.Print(out)
	}
	return
}
