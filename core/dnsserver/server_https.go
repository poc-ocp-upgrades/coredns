package dnsserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
	"github.com/coredns/coredns/plugin/pkg/dnsutil"
	"github.com/coredns/coredns/plugin/pkg/doh"
	"github.com/coredns/coredns/plugin/pkg/response"
	"github.com/coredns/coredns/plugin/pkg/transport"
)

type ServerHTTPS struct {
	*Server
	httpsServer	*http.Server
	listenAddr	net.Addr
	tlsConfig	*tls.Config
}

func NewServerHTTPS(addr string, group []*Config) (*ServerHTTPS, error) {
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
	sh := &ServerHTTPS{Server: s, tlsConfig: tlsConfig, httpsServer: new(http.Server)}
	sh.httpsServer.Handler = sh
	return sh, nil
}
func (s *ServerHTTPS) Serve(l net.Listener) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.m.Lock()
	s.listenAddr = l.Addr()
	s.m.Unlock()
	if s.tlsConfig != nil {
		l = tls.NewListener(l, s.tlsConfig)
	}
	return s.httpsServer.Serve(l)
}
func (s *ServerHTTPS) ServePacket(p net.PacketConn) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (s *ServerHTTPS) Listen() (net.Listener, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l, err := net.Listen("tcp", s.Addr[len(transport.HTTPS+"://"):])
	if err != nil {
		return nil, err
	}
	return l, nil
}
func (s *ServerHTTPS) ListenPacket() (net.PacketConn, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, nil
}
func (s *ServerHTTPS) OnStartupComplete() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if Quiet {
		return
	}
	out := startUpZones(transport.HTTPS+"://", s.Addr, s.zones)
	if out != "" {
		fmt.Print(out)
	}
	return
}
func (s *ServerHTTPS) Stop() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.m.Lock()
	defer s.m.Unlock()
	if s.httpsServer != nil {
		s.httpsServer.Shutdown(context.Background())
	}
	return nil
}
func (s *ServerHTTPS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.URL.Path != doh.Path {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	msg, err := doh.RequestToMsg(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h, p, _ := net.SplitHostPort(r.RemoteAddr)
	port, _ := strconv.Atoi(p)
	dw := &DoHWriter{laddr: s.listenAddr, raddr: &net.TCPAddr{IP: net.ParseIP(h), Port: port}}
	s.ServeDNS(context.Background(), dw, msg)
	buf, _ := dw.Msg.Pack()
	mt, _ := response.Typify(dw.Msg, time.Now().UTC())
	age := dnsutil.MinimalTTL(dw.Msg, mt)
	w.Header().Set("Content-Type", doh.MimeType)
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%f", age.Seconds()))
	w.Header().Set("Content-Length", strconv.Itoa(len(buf)))
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}
func (s *ServerHTTPS) Shutdown() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.httpsServer != nil {
		s.httpsServer.Shutdown(context.Background())
	}
	return nil
}
