package dnsserver

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"sync"
	"time"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics/vars"
	"github.com/coredns/coredns/plugin/pkg/edns"
	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/plugin/pkg/rcode"
	"github.com/coredns/coredns/plugin/pkg/trace"
	"github.com/coredns/coredns/plugin/pkg/transport"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	ot "github.com/opentracing/opentracing-go"
)

type Server struct {
	Addr		string
	server		[2]*dns.Server
	m		sync.Mutex
	zones		map[string]*Config
	dnsWg		sync.WaitGroup
	connTimeout	time.Duration
	trace		trace.Trace
	debug		bool
	classChaos	bool
}

func NewServer(addr string, group []*Config) (*Server, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := &Server{Addr: addr, zones: make(map[string]*Config), connTimeout: 5 * time.Second}
	s.dnsWg.Add(1)
	for _, site := range group {
		if site.Debug {
			s.debug = true
			log.D = true
		}
		s.zones[site.Zone] = site
		if site.registry != nil {
			for name := range EnableChaos {
				if _, ok := site.registry[name]; ok {
					s.classChaos = true
					break
				}
			}
			if handler, ok := site.registry["trace"]; ok {
				s.trace = handler.(trace.Trace)
			}
			continue
		}
		var stack plugin.Handler
		for i := len(site.Plugin) - 1; i >= 0; i-- {
			stack = site.Plugin[i](stack)
			site.registerHandler(stack)
			if s.trace == nil && stack.Name() == "trace" {
				if t, ok := stack.(trace.Trace); ok {
					s.trace = t
				}
			}
			if _, ok := EnableChaos[stack.Name()]; ok {
				s.classChaos = true
			}
		}
		site.pluginChain = stack
	}
	return s, nil
}
func (s *Server) Serve(l net.Listener) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.m.Lock()
	s.server[tcp] = &dns.Server{Listener: l, Net: "tcp", Handler: dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
		ctx := context.WithValue(context.Background(), Key{}, s)
		s.ServeDNS(ctx, w, r)
	})}
	s.m.Unlock()
	return s.server[tcp].ActivateAndServe()
}
func (s *Server) ServePacket(p net.PacketConn) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.m.Lock()
	s.server[udp] = &dns.Server{PacketConn: p, Net: "udp", Handler: dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
		ctx := context.WithValue(context.Background(), Key{}, s)
		s.ServeDNS(ctx, w, r)
	})}
	s.m.Unlock()
	return s.server[udp].ActivateAndServe()
}
func (s *Server) Listen() (net.Listener, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l, err := listen("tcp", s.Addr[len(transport.DNS+"://"):])
	if err != nil {
		return nil, err
	}
	return l, nil
}
func (s *Server) ListenPacket() (net.PacketConn, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	p, err := listenPacket("udp", s.Addr[len(transport.DNS+"://"):])
	if err != nil {
		return nil, err
	}
	return p, nil
}
func (s *Server) Stop() (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if runtime.GOOS != "windows" {
		done := make(chan struct{})
		go func() {
			s.dnsWg.Done()
			s.dnsWg.Wait()
			close(done)
		}()
		select {
		case <-time.After(s.connTimeout):
		case <-done:
		}
	}
	s.m.Lock()
	for _, s1 := range s.server {
		if s1 != nil {
			err = s1.Shutdown()
		}
	}
	s.m.Unlock()
	return
}
func (s *Server) Address() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return s.Addr
}
func (s *Server) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r == nil || len(r.Question) == 0 {
		DefaultErrorFunc(ctx, w, r, dns.RcodeServerFailure)
		return
	}
	if !s.debug {
		defer func() {
			if rec := recover(); rec != nil {
				vars.Panic.Inc()
				DefaultErrorFunc(ctx, w, r, dns.RcodeServerFailure)
			}
		}()
	}
	if !s.classChaos && r.Question[0].Qclass != dns.ClassINET {
		DefaultErrorFunc(ctx, w, r, dns.RcodeRefused)
		return
	}
	if m, err := edns.Version(r); err != nil {
		w.WriteMsg(m)
		return
	}
	ctx, err := incrementDepthAndCheck(ctx)
	if err != nil {
		DefaultErrorFunc(ctx, w, r, dns.RcodeServerFailure)
		return
	}
	q := r.Question[0].Name
	b := make([]byte, len(q))
	var off int
	var end bool
	var dshandler *Config
	w = request.NewScrubWriter(r, w)
	for {
		l := len(q[off:])
		for i := 0; i < l; i++ {
			b[i] = q[off+i]
			if b[i] >= 'A' && b[i] <= 'Z' {
				b[i] |= ('a' - 'A')
			}
		}
		if h, ok := s.zones[string(b[:l])]; ok {
			ctx = context.WithValue(ctx, plugin.ServerCtx{}, s.Addr)
			if r.Question[0].Qtype != dns.TypeDS {
				if h.FilterFunc == nil {
					rcode, _ := h.pluginChain.ServeDNS(ctx, w, r)
					if !plugin.ClientWrite(rcode) {
						DefaultErrorFunc(ctx, w, r, rcode)
					}
					return
				}
				if h.FilterFunc(q) {
					rcode, _ := h.pluginChain.ServeDNS(ctx, w, r)
					if !plugin.ClientWrite(rcode) {
						DefaultErrorFunc(ctx, w, r, rcode)
					}
					return
				}
			}
			dshandler = h
		}
		off, end = dns.NextLabel(q, off)
		if end {
			break
		}
	}
	if r.Question[0].Qtype == dns.TypeDS && dshandler != nil && dshandler.pluginChain != nil {
		rcode, _ := dshandler.pluginChain.ServeDNS(ctx, w, r)
		if !plugin.ClientWrite(rcode) {
			DefaultErrorFunc(ctx, w, r, rcode)
		}
		return
	}
	if h, ok := s.zones["."]; ok && h.pluginChain != nil {
		ctx = context.WithValue(ctx, plugin.ServerCtx{}, s.Addr)
		rcode, _ := h.pluginChain.ServeDNS(ctx, w, r)
		if !plugin.ClientWrite(rcode) {
			DefaultErrorFunc(ctx, w, r, rcode)
		}
		return
	}
	DefaultErrorFunc(ctx, w, r, dns.RcodeRefused)
}
func (s *Server) OnStartupComplete() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if Quiet {
		return
	}
	out := startUpZones("", s.Addr, s.zones)
	if out != "" {
		fmt.Print(out)
	}
	return
}
func (s *Server) Tracer() ot.Tracer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.trace == nil {
		return nil
	}
	return s.trace.Tracer()
}
func DefaultErrorFunc(ctx context.Context, w dns.ResponseWriter, r *dns.Msg, rc int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	state := request.Request{W: w, Req: r}
	answer := new(dns.Msg)
	answer.SetRcode(r, rc)
	state.SizeAndDo(answer)
	vars.Report(ctx, state, vars.Dropped, rcode.ToString(rc), answer.Len(), time.Now())
	w.WriteMsg(answer)
}
func incrementDepthAndCheck(ctx context.Context) (context.Context, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	loop := ctx.Value(loopKey{})
	if loop == nil {
		ctx = context.WithValue(ctx, loopKey{}, 0)
		return ctx, nil
	}
	iloop := loop.(int) + 1
	if iloop > maxreentries {
		return ctx, fmt.Errorf("too deep")
	}
	ctx = context.WithValue(ctx, loopKey{}, iloop)
	return ctx, nil
}

const (
	tcp		= 0
	udp		= 1
	maxreentries	= 10
)

type (
	Key	struct{}
	loopKey	struct{}
)

var EnableChaos = map[string]struct{}{"chaos": struct{}{}, "forward": struct{}{}, "proxy": struct{}{}}
var Quiet bool
