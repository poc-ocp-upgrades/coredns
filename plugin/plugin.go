package plugin

import (
	"context"
	"errors"
	"fmt"
	"github.com/miekg/dns"
	ot "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
)

type (
	Plugin	func(Handler) Handler
	Handler	interface {
		ServeDNS(context.Context, dns.ResponseWriter, *dns.Msg) (int, error)
		Name() string
	}
	HandlerFunc	func(context.Context, dns.ResponseWriter, *dns.Msg) (int, error)
)

func (f HandlerFunc) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return f(ctx, w, r)
}
func (f HandlerFunc) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "handlerfunc"
}
func Error(name string, err error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Errorf("%s/%s: %s", "plugin", name, err)
}
func NextOrFailure(name string, next Handler, ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if next != nil {
		if span := ot.SpanFromContext(ctx); span != nil {
			child := span.Tracer().StartSpan(next.Name(), ot.ChildOf(span.Context()))
			defer child.Finish()
			ctx = ot.ContextWithSpan(ctx, child)
		}
		return next.ServeDNS(ctx, w, r)
	}
	return dns.RcodeServerFailure, Error(name, errors.New("no next plugin found"))
}
func ClientWrite(rcode int) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch rcode {
	case dns.RcodeServerFailure:
		fallthrough
	case dns.RcodeRefused:
		fallthrough
	case dns.RcodeFormatError:
		fallthrough
	case dns.RcodeNotImplemented:
		return false
	}
	return true
}

const Namespace = "coredns"

var TimeBuckets = prometheus.ExponentialBuckets(0.00025, 2, 16)
var ErrOnce = errors.New("this plugin can only be used once per Server Block")

type ServerCtx struct{}
