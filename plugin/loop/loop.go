package loop

import (
	"context"
	"sync"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

var log = clog.NewWithPlugin("loop")

type Loop struct {
	Next	plugin.Handler
	zone	string
	qname	string
	addr	string
	sync.RWMutex
	i	int
	off	bool
}

func New(zone string) *Loop {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &Loop{zone: zone, qname: qname(zone)}
}
func (l *Loop) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.Question[0].Qtype != dns.TypeHINFO {
		return plugin.NextOrFailure(l.Name(), l.Next, ctx, w, r)
	}
	if l.disabled() {
		return plugin.NextOrFailure(l.Name(), l.Next, ctx, w, r)
	}
	state := request.Request{W: w, Req: r}
	zone := plugin.Zones([]string{l.zone}).Matches(state.Name())
	if zone == "" {
		return plugin.NextOrFailure(l.Name(), l.Next, ctx, w, r)
	}
	if state.Name() == l.qname {
		l.inc()
	}
	if l.seen() > 2 {
		log.Fatalf(`Loop (%s -> %s) detected for zone %q, see https://coredns.io/plugins/loop#troubleshooting. Query: "HINFO %s"`, state.RemoteAddr(), l.address(), l.zone, l.qname)
	}
	return plugin.NextOrFailure(l.Name(), l.Next, ctx, w, r)
}
func (l *Loop) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "loop"
}
func (l *Loop) exchange(addr string) (*dns.Msg, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := new(dns.Msg)
	m.SetQuestion(l.qname, dns.TypeHINFO)
	return dns.Exchange(m, addr)
}
func (l *Loop) seen() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l.RLock()
	defer l.RUnlock()
	return l.i
}
func (l *Loop) inc() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l.Lock()
	defer l.Unlock()
	l.i++
}
func (l *Loop) reset() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l.Lock()
	defer l.Unlock()
	l.i = 0
}
func (l *Loop) setDisabled() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l.Lock()
	defer l.Unlock()
	l.off = true
}
func (l *Loop) disabled() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l.RLock()
	defer l.RUnlock()
	return l.off
}
func (l *Loop) setAddress(addr string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l.Lock()
	defer l.Unlock()
	l.addr = addr
}
func (l *Loop) address() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l.RLock()
	defer l.RUnlock()
	return l.addr
}
