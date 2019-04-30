package kubernetes

import (
	"strings"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/plugin/pkg/dnsutil"
	"github.com/coredns/coredns/request"
)

func (k *Kubernetes) Reverse(state request.Request, exact bool, opt plugin.Options) ([]msg.Service, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ip := dnsutil.ExtractAddressFromReverse(state.Name())
	if ip == "" {
		_, e := k.Records(state, exact)
		return nil, e
	}
	records := k.serviceRecordForIP(ip, state.Name())
	if len(records) == 0 {
		return records, errNoItems
	}
	return records, nil
}
func (k *Kubernetes) serviceRecordForIP(ip, name string) []msg.Service {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, service := range k.APIConn.SvcIndexReverse(ip) {
		if len(k.Namespaces) > 0 && !k.namespaceExposed(service.Namespace) {
			continue
		}
		domain := strings.Join([]string{service.Name, service.Namespace, Svc, k.primaryZone()}, ".")
		return []msg.Service{{Host: domain, TTL: k.ttl}}
	}
	for _, ep := range k.APIConn.EpIndexReverse(ip) {
		if len(k.Namespaces) > 0 && !k.namespaceExposed(ep.Namespace) {
			continue
		}
		for _, eps := range ep.Subsets {
			for _, addr := range eps.Addresses {
				if addr.IP == ip {
					domain := strings.Join([]string{endpointHostname(addr, k.endpointNameMode), ep.Name, ep.Namespace, Svc, k.primaryZone()}, ".")
					return []msg.Service{{Host: domain, TTL: k.ttl}}
				}
			}
		}
	}
	return nil
}
