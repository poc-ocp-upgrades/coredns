package kubernetes

import (
	"net"
)

func localPodIP() net.IP {
	_logClusterCodePath()
	defer _logClusterCodePath()
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, addr := range addrs {
		ip, _, _ := net.ParseCIDR(addr.String())
		ip = ip.To4()
		if ip == nil || ip.IsLoopback() {
			continue
		}
		return ip
	}
	return nil
}
func (k *Kubernetes) localNodeName() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	localIP := k.interfaceAddrsFunc()
	if localIP == nil {
		return ""
	}
	for _, ep := range k.APIConn.EpIndexReverse(localIP.String()) {
		for _, eps := range ep.Subsets {
			for _, addr := range eps.Addresses {
				if localIP.Equal(net.ParseIP(addr.IP)) {
					return addr.NodeName
				}
			}
		}
	}
	return ""
}
