package dnsserver

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"github.com/coredns/coredns/pb"
	"github.com/coredns/coredns/plugin/pkg/transport"
	"github.com/coredns/coredns/plugin/pkg/watch"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/miekg/dns"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type ServergRPC struct {
	*Server
	grpcServer	*grpc.Server
	listenAddr	net.Addr
	tlsConfig	*tls.Config
	watch		watch.Watcher
}

func NewServergRPC(addr string, group []*Config) (*ServergRPC, error) {
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
	return &ServergRPC{Server: s, tlsConfig: tlsConfig, watch: watch.NewWatcher(watchables(s.zones))}, nil
}
func (s *ServergRPC) Serve(l net.Listener) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.m.Lock()
	s.listenAddr = l.Addr()
	s.m.Unlock()
	if s.Tracer() != nil {
		onlyIfParent := func(parentSpanCtx opentracing.SpanContext, method string, req, resp interface{}) bool {
			return parentSpanCtx != nil
		}
		intercept := otgrpc.OpenTracingServerInterceptor(s.Tracer(), otgrpc.IncludingSpans(onlyIfParent))
		s.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(intercept))
	} else {
		s.grpcServer = grpc.NewServer()
	}
	pb.RegisterDnsServiceServer(s.grpcServer, s)
	if s.tlsConfig != nil {
		l = tls.NewListener(l, s.tlsConfig)
	}
	return s.grpcServer.Serve(l)
}
func (s *ServergRPC) ServePacket(p net.PacketConn) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (s *ServergRPC) Listen() (net.Listener, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	l, err := net.Listen("tcp", s.Addr[len(transport.GRPC+"://"):])
	if err != nil {
		return nil, err
	}
	return l, nil
}
func (s *ServergRPC) ListenPacket() (net.PacketConn, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, nil
}
func (s *ServergRPC) OnStartupComplete() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if Quiet {
		return
	}
	out := startUpZones(transport.GRPC+"://", s.Addr, s.zones)
	if out != "" {
		fmt.Print(out)
	}
	return
}
func (s *ServergRPC) Stop() (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.m.Lock()
	defer s.m.Unlock()
	if s.watch != nil {
		s.watch.Stop()
	}
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
	return
}
func (s *ServergRPC) Query(ctx context.Context, in *pb.DnsPacket) (*pb.DnsPacket, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	msg := new(dns.Msg)
	err := msg.Unpack(in.Msg)
	if err != nil {
		return nil, err
	}
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("no peer in gRPC context")
	}
	a, ok := p.Addr.(*net.TCPAddr)
	if !ok {
		return nil, fmt.Errorf("no TCP peer in gRPC context: %v", p.Addr)
	}
	w := &gRPCresponse{localAddr: s.listenAddr, remoteAddr: a, Msg: msg}
	s.ServeDNS(ctx, w, msg)
	packed, err := w.Msg.Pack()
	if err != nil {
		return nil, err
	}
	return &pb.DnsPacket{Msg: packed}, nil
}
func (s *ServergRPC) Watch(stream pb.DnsService_WatchServer) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return s.watch.Watch(stream)
}
func (s *ServergRPC) Shutdown() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}
	return nil
}

type gRPCresponse struct {
	localAddr	net.Addr
	remoteAddr	net.Addr
	Msg		*dns.Msg
}

func (r *gRPCresponse) Write(b []byte) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.Msg = new(dns.Msg)
	return len(b), r.Msg.Unpack(b)
}
func (r *gRPCresponse) Close() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (r *gRPCresponse) TsigStatus() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (r *gRPCresponse) TsigTimersOnly(b bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return
}
func (r *gRPCresponse) Hijack() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return
}
func (r *gRPCresponse) LocalAddr() net.Addr {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.localAddr
}
func (r *gRPCresponse) RemoteAddr() net.Addr {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.remoteAddr
}
func (r *gRPCresponse) WriteMsg(m *dns.Msg) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.Msg = m
	return nil
}
