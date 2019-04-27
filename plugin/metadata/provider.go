package metadata

import (
	"context"
	"strings"
	"github.com/coredns/coredns/request"
)

type Provider interface {
	Metadata(ctx context.Context, state request.Request) context.Context
}
type Func func() string

func IsLabel(label string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	p := strings.Index(label, "/")
	if p <= 0 || p >= len(label)-1 {
		return false
	}
	if strings.LastIndex(label, "/") != p {
		return false
	}
	return true
}
func Labels(ctx context.Context) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if metadata := ctx.Value(key{}); metadata != nil {
		if m, ok := metadata.(md); ok {
			return keys(m)
		}
	}
	return nil
}
func ValueFunc(ctx context.Context, label string) Func {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if metadata := ctx.Value(key{}); metadata != nil {
		if m, ok := metadata.(md); ok {
			return m[label]
		}
	}
	return nil
}
func SetValueFunc(ctx context.Context, label string, f Func) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if metadata := ctx.Value(key{}); metadata != nil {
		if m, ok := metadata.(md); ok {
			m[label] = f
			return true
		}
	}
	return false
}

type md map[string]Func
type key struct{}

func keys(m map[string]Func) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := make([]string, len(m))
	i := 0
	for k := range m {
		s[i] = k
		i++
	}
	return s
}
