package kubernetes

import (
	"github.com/coredns/coredns/plugin/kubernetes/object"
	"github.com/coredns/coredns/plugin/pkg/watch"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kubernetes) SetWatchChan(c watch.Chan) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	k.APIConn.SetWatchChan(c)
}
func (k *Kubernetes) Watch(qname string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return k.APIConn.Watch(qname)
}
func (k *Kubernetes) StopWatching(qname string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	k.APIConn.StopWatching(qname)
}

var _ watch.Watchable = &Kubernetes{}

func (dns *dnsControl) sendServiceUpdates(s *object.Service) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for i := range dns.zones {
		name := serviceFQDN(s, dns.zones[i])
		if _, ok := dns.watched[name]; ok {
			dns.watchChan <- name
		}
	}
}
func (dns *dnsControl) sendPodUpdates(p *object.Pod) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for i := range dns.zones {
		name := podFQDN(p, dns.zones[i])
		if _, ok := dns.watched[name]; ok {
			dns.watchChan <- name
		}
	}
}
func (dns *dnsControl) sendEndpointsUpdates(ep *object.Endpoints) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, zone := range dns.zones {
		for _, name := range endpointFQDN(ep, zone, dns.endpointNameMode) {
			if _, ok := dns.watched[name]; ok {
				dns.watchChan <- name
			}
		}
		name := serviceFQDN(ep, zone)
		if _, ok := dns.watched[name]; ok {
			dns.watchChan <- name
		}
	}
}
func endpointsSubsetDiffs(a, b *object.Endpoints) *object.Endpoints {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := b.CopyWithoutSubsets()
	for _, abba := range [][]*object.Endpoints{{a, b}, {b, a}} {
		a := abba[0]
		b := abba[1]
	left:
		for _, as := range a.Subsets {
			for _, bs := range b.Subsets {
				if subsetsEquivalent(as, bs) {
					continue left
				}
			}
			c.Subsets = append(c.Subsets, as)
		}
	}
	return c
}
func (dns *dnsControl) sendUpdates(oldObj, newObj interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if newObj != nil && oldObj != nil && (oldObj.(meta.Object).GetResourceVersion() == newObj.(meta.Object).GetResourceVersion()) {
		return
	}
	obj := newObj
	if obj == nil {
		obj = oldObj
	}
	switch ob := obj.(type) {
	case *object.Service:
		dns.updateModifed()
		dns.sendServiceUpdates(ob)
	case *object.Endpoints:
		if newObj == nil || oldObj == nil {
			dns.updateModifed()
			dns.sendEndpointsUpdates(ob)
			return
		}
		p := oldObj.(*object.Endpoints)
		if endpointsEquivalent(p, ob) {
			return
		}
		dns.updateModifed()
		dns.sendEndpointsUpdates(endpointsSubsetDiffs(p, ob))
	case *object.Pod:
		dns.updateModifed()
		dns.sendPodUpdates(ob)
	default:
		log.Warningf("Updates for %T not supported.", ob)
	}
}
func (dns *dnsControl) Add(obj interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	dns.sendUpdates(nil, obj)
}
func (dns *dnsControl) Delete(obj interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	dns.sendUpdates(obj, nil)
}
func (dns *dnsControl) Update(oldObj, newObj interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	dns.sendUpdates(oldObj, newObj)
}
func subsetsEquivalent(sa, sb object.EndpointSubset) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(sa.Addresses) != len(sb.Addresses) {
		return false
	}
	if len(sa.Ports) != len(sb.Ports) {
		return false
	}
	for addr, aaddr := range sa.Addresses {
		baddr := sb.Addresses[addr]
		if aaddr.IP != baddr.IP {
			return false
		}
		if aaddr.Hostname != baddr.Hostname {
			return false
		}
	}
	for port, aport := range sa.Ports {
		bport := sb.Ports[port]
		if aport.Name != bport.Name {
			return false
		}
		if aport.Port != bport.Port {
			return false
		}
		if aport.Protocol != bport.Protocol {
			return false
		}
	}
	return true
}
func endpointsEquivalent(a, b *object.Endpoints) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(a.Subsets) != len(b.Subsets) {
		return false
	}
	for i, sa := range a.Subsets {
		sb := b.Subsets[i]
		if !subsetsEquivalent(sa, sb) {
			return false
		}
	}
	return true
}
