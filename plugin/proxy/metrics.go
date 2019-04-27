package proxy

import (
	"github.com/coredns/coredns/plugin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestCount	= prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: plugin.Namespace, Subsystem: "proxy", Name: "request_count_total", Help: "Counter of requests made per protocol, proxy protocol, family and upstream."}, []string{"server", "proto", "proxy_proto", "family", "to"})
	RequestDuration	= prometheus.NewHistogramVec(prometheus.HistogramOpts{Namespace: plugin.Namespace, Subsystem: "proxy", Name: "request_duration_seconds", Buckets: plugin.TimeBuckets, Help: "Histogram of the time (in seconds) each request took."}, []string{"server", "proto", "proxy_proto", "family", "to"})
)

func familyToString(f int) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if f == 1 {
		return "1"
	}
	if f == 2 {
		return "2"
	}
	return ""
}
