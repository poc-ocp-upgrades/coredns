package kubernetes

import (
	"strings"
	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/plugin/kubernetes/object"
	"github.com/coredns/coredns/plugin/pkg/dnsutil"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func (k *Kubernetes) External(state request.Request) ([]msg.Service, int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	base, _ := dnsutil.TrimZone(state.Name(), state.Zone)
	segs := dns.SplitDomainName(base)
	last := len(segs) - 1
	if last < 0 {
		return nil, dns.RcodeServerFailure
	}
	port := "*"
	protocol := "*"
	namespace := segs[last]
	if !k.namespaceExposed(namespace) || !k.namespace(namespace) {
		return nil, dns.RcodeNameError
	}
	last--
	if last < 0 {
		return nil, dns.RcodeSuccess
	}
	service := segs[last]
	last--
	if last == 1 {
		protocol = stripUnderscore(segs[last])
		port = stripUnderscore(segs[last-1])
		last -= 2
	}
	if last != -1 {
		return nil, dns.RcodeNameError
	}
	idx := object.ServiceKey(service, namespace)
	serviceList := k.APIConn.SvcIndex(idx)
	services := []msg.Service{}
	zonePath := msg.Path(state.Zone, coredns)
	rcode := dns.RcodeNameError
	for _, svc := range serviceList {
		if namespace != svc.Namespace {
			continue
		}
		if service != svc.Name {
			continue
		}
		for _, ip := range svc.ExternalIPs {
			for _, p := range svc.Ports {
				if !(match(port, p.Name) && match(protocol, string(p.Protocol))) {
					continue
				}
				rcode = dns.RcodeSuccess
				s := msg.Service{Host: ip, Port: int(p.Port), TTL: k.ttl}
				s.Key = strings.Join([]string{zonePath, svc.Namespace, svc.Name}, "/")
				services = append(services, s)
			}
		}
	}
	return services, rcode
}
func (k *Kubernetes) ExternalAddress(state request.Request) []dns.RR {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rrs := []dns.RR{k.nsAddr()}
	return rrs
}
