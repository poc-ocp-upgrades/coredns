package etcd

import (
	"context"
	"time"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func (e *Etcd) Serial(state request.Request) uint32 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return uint32(time.Now().Unix())
}
func (e *Etcd) MinTTL(state request.Request) uint32 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return 30
}
func (e *Etcd) Transfer(ctx context.Context, state request.Request) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return dns.RcodeServerFailure, nil
}
