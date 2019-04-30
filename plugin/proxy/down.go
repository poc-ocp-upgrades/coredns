package proxy

import (
	"sync/atomic"
	"github.com/coredns/coredns/plugin/pkg/healthcheck"
)

var checkDownFunc = func(upstream *staticUpstream) healthcheck.UpstreamHostDownFunc {
	return func(uh *healthcheck.UpstreamHost) bool {
		fails := atomic.LoadInt32(&uh.Fails)
		if fails >= upstream.MaxFails && upstream.MaxFails != 0 {
			return true
		}
		return false
	}
}
