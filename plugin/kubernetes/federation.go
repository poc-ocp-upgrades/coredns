package kubernetes

import (
	"errors"
	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/plugin/pkg/dnsutil"
	"github.com/coredns/coredns/request"
)

const (
	LabelZone	= "failure-domain.beta.kubernetes.io/zone"
	LabelRegion	= "failure-domain.beta.kubernetes.io/region"
)

func (k *Kubernetes) Federations(state request.Request, fname, fzone string) (msg.Service, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	nodeName := k.localNodeName()
	node, err := k.APIConn.GetNodeByName(nodeName)
	if err != nil {
		return msg.Service{}, err
	}
	r, err := parseRequest(state)
	if err != nil {
		return msg.Service{}, err
	}
	lz := node.Labels[LabelZone]
	lr := node.Labels[LabelRegion]
	if lz == "" || lr == "" {
		return msg.Service{}, errors.New("local node missing zone/region labels")
	}
	if r.endpoint == "" {
		return msg.Service{Host: dnsutil.Join(r.service, r.namespace, fname, r.podOrSvc, lz, lr, fzone)}, nil
	}
	return msg.Service{Host: dnsutil.Join(r.endpoint, r.service, r.namespace, fname, r.podOrSvc, lz, lr, fzone)}, nil
}
