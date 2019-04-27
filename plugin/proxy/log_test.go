package proxy

import clog "github.com/coredns/coredns/plugin/pkg/log"

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	clog.Discard()
}
