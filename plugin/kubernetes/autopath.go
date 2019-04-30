package kubernetes

import (
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/kubernetes/object"
	"github.com/coredns/coredns/request"
)

func (k *Kubernetes) AutoPath(state request.Request) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	zone := plugin.Zones(k.Zones).Matches(state.Name())
	if zone == "" {
		return nil
	}
	if !k.opts.initPodCache {
		return nil
	}
	ip := state.IP()
	pod := k.podWithIP(ip)
	if pod == nil {
		return nil
	}
	search := make([]string, 3)
	if zone == "." {
		search[0] = pod.Namespace + ".svc."
		search[1] = "svc."
		search[2] = "."
	} else {
		search[0] = pod.Namespace + ".svc." + zone
		search[1] = "svc." + zone
		search[2] = zone
	}
	search = append(search, k.autoPathSearch...)
	search = append(search, "")
	return search
}
func (k *Kubernetes) podWithIP(ip string) *object.Pod {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ps := k.APIConn.PodIndex(ip)
	if len(ps) == 0 {
		return nil
	}
	return ps[0]
}
