package proxy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"
	"github.com/mholt/caddy/caddyfile"
	"github.com/miekg/dns"
)

func TestStop(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	config := "proxy . %s {\n health_check /healthcheck:%s %dms \n}"
	tests := []struct {
		intervalInMilliseconds	int
		numHealthcheckIntervals	int
	}{{5, 1}, {5, 2}, {5, 3}}
	for i, test := range tests {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			var counter int64
			backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.Body.Close()
				atomic.AddInt64(&counter, 1)
			}))
			defer backend.Close()
			port := backend.URL[17:]
			back := backend.URL[7:]
			c := caddyfile.NewDispenser("Testfile", strings.NewReader(fmt.Sprintf(config, back, port, test.intervalInMilliseconds)))
			upstreams, err := NewStaticUpstreams(&c)
			if err != nil {
				t.Errorf("Test %d, expected no error. Got: %s", i, err)
			}
			time.Sleep(time.Duration(test.intervalInMilliseconds*test.numHealthcheckIntervals) * time.Millisecond)
			for _, upstream := range upstreams {
				if err := upstream.Stop(); err != nil {
					t.Errorf("Test %d, expected no error stopping upstream, got: %s", i, err)
				}
			}
			counterAfterShutdown := atomic.LoadInt64(&counter)
			time.Sleep(time.Duration(test.intervalInMilliseconds*test.numHealthcheckIntervals) * time.Millisecond)
			if counterAfterShutdown == 0 {
				t.Errorf("Test %d, Expected healthchecks to hit test server, got none", i)
			}
			counterAfterWaiting := atomic.LoadInt64(&counter)
			if counterAfterWaiting > (counterAfterShutdown + 1) {
				t.Errorf("Test %d, expected no more healthchecks after shutdown. got: %d healthchecks after shutdown", i, counterAfterWaiting-counterAfterShutdown)
			}
		})
	}
}
func TestProxyRefused(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := dnstest.NewServer(func(w dns.ResponseWriter, r *dns.Msg) {
		ret := new(dns.Msg)
		ret.SetReply(r)
		ret.Rcode = dns.RcodeRefused
		w.WriteMsg(ret)
	})
	defer s.Close()
	p := NewLookup([]string{s.Addr})
	state := request.Request{W: &test.ResponseWriter{}, Req: new(dns.Msg)}
	state.Req.SetQuestion("example.org.", dns.TypeA)
	resp, err := p.Forward(state)
	if err != nil {
		t.Fatal("Expected to receive reply, but didn't")
	}
	if resp.Rcode != dns.RcodeRefused {
		t.Errorf("Expected rcode to be %d, got %d", dns.RcodeRefused, resp.Rcode)
	}
}
