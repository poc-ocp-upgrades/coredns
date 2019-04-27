package dnsserver

import (
	"github.com/coredns/coredns/plugin/pkg/watch"
)

func watchables(zones map[string]*Config) []watch.Watchable {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var w []watch.Watchable
	for _, config := range zones {
		plugins := config.Handlers()
		for _, p := range plugins {
			if x, ok := p.(watch.Watchable); ok {
				w = append(w, x)
			}
		}
	}
	return w
}
