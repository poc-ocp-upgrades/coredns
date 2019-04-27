package proxy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
	"github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"
	"github.com/mholt/caddy/caddyfile"
	"github.com/miekg/dns"
)

func TestUnhealthy(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	config := "proxy . %s {\n health_check /healthcheck:%s 10s \nfail_timeout 100ms\n}"
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body.Close()
		w.Write([]byte("OK"))
	}))
	defer backend.Close()
	port := backend.URL[17:]
	back := backend.URL[7:]
	c := caddyfile.NewDispenser("testfile", strings.NewReader(fmt.Sprintf(config, back, port)))
	upstreams, err := NewStaticUpstreams(&c)
	if err != nil {
		t.Errorf("Expected no error. Got: %s", err)
	}
	p := &Proxy{Upstreams: &upstreams}
	m := new(dns.Msg)
	m.SetQuestion("example.org.", dns.TypeA)
	state := request.Request{W: &test.ResponseWriter{}, Req: m}
	for j := 0; j < failureCheck; j++ {
		if _, err := p.Forward(state); err == nil {
			t.Errorf("Expected error. Got: nil")
		}
	}
	fails := atomic.LoadInt32(&upstreams[0].(*staticUpstream).Hosts[0].Fails)
	if fails != 3 {
		t.Errorf("Expected %d fails, got %d", 3, fails)
	}
	i := 0
	for fails != 0 {
		fails = atomic.LoadInt32(&upstreams[0].(*staticUpstream).Hosts[0].Fails)
		time.Sleep(100 * time.Microsecond)
		i++
	}
}
