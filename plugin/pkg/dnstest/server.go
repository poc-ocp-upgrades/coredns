package dnstest

import (
	"net"
	"github.com/miekg/dns"
)

type Server struct {
	Addr	string
	s1		*dns.Server
	s2		*dns.Server
}

func NewServer(f dns.HandlerFunc) *Server {
	_logClusterCodePath()
	defer _logClusterCodePath()
	dns.HandleFunc(".", f)
	ch1 := make(chan bool)
	ch2 := make(chan bool)
	s1 := &dns.Server{}
	s2 := &dns.Server{}
	for i := 0; i < 5; i++ {
		s2.Listener, _ = net.Listen("tcp", ":0")
		if s2.Listener == nil {
			continue
		}
		s1.PacketConn, _ = net.ListenPacket("udp", s2.Listener.Addr().String())
		if s1.PacketConn != nil {
			break
		}
		s2.Listener.Close()
		s2.Listener = nil
	}
	if s2.Listener == nil {
		panic("dnstest.NewServer(): failed to create new server")
	}
	s1.NotifyStartedFunc = func() {
		close(ch1)
	}
	s2.NotifyStartedFunc = func() {
		close(ch2)
	}
	go s1.ActivateAndServe()
	go s2.ActivateAndServe()
	<-ch1
	<-ch2
	return &Server{s1: s1, s2: s2, Addr: s2.Listener.Addr().String()}
}
func (s *Server) Close() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.s1.Shutdown()
	s.s2.Shutdown()
}
