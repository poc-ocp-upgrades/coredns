package template

import (
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	"github.com/mholt/caddy"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	templateMatchesCount	= prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: plugin.Namespace, Subsystem: "template", Name: "matches_total", Help: "Counter of template regex matches."}, []string{"server", "zone", "class", "type"})
	templateFailureCount	= prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: plugin.Namespace, Subsystem: "template", Name: "template_failures_total", Help: "Counter of go template failures."}, []string{"server", "zone", "class", "type", "section", "template"})
	templateRRFailureCount	= prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: plugin.Namespace, Subsystem: "template", Name: "rr_failures_total", Help: "Counter of mis-templated RRs."}, []string{"server", "zone", "class", "type", "section", "template"})
)

func setupMetrics(c *caddy.Controller) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.OnStartup(func() error {
		metrics.MustRegister(c, templateMatchesCount, templateFailureCount, templateRRFailureCount)
		return nil
	})
	return nil
}
