package federation

import (
	"github.com/coredns/coredns/plugin/kubernetes"
	"github.com/coredns/coredns/plugin/kubernetes/object"
	"github.com/coredns/coredns/plugin/pkg/watch"
	api "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type APIConnFederationTest struct{ zone, region string }

func (APIConnFederationTest) HasSynced() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return true
}
func (APIConnFederationTest) Run() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return
}
func (APIConnFederationTest) Stop() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (APIConnFederationTest) SvcIndexReverse(string) []*object.Service {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (APIConnFederationTest) EpIndexReverse(string) []*object.Endpoints {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (APIConnFederationTest) Modified() int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return 0
}
func (APIConnFederationTest) SetWatchChan(watch.Chan) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (APIConnFederationTest) Watch(string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (APIConnFederationTest) StopWatching(string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (APIConnFederationTest) PodIndex(string) []*object.Pod {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return []*object.Pod{{Namespace: "podns", PodIP: "10.240.0.1"}}
}
func (APIConnFederationTest) SvcIndex(string) []*object.Service {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	svcs := []*object.Service{{Name: "svc1", Namespace: "testns", ClusterIP: "10.0.0.1", Ports: []api.ServicePort{{Name: "http", Protocol: "tcp", Port: 80}}}, {Name: "hdls1", Namespace: "testns", ClusterIP: api.ClusterIPNone}, {Name: "external", Namespace: "testns", ExternalName: "ext.interwebs.test", Ports: []api.ServicePort{{Name: "http", Protocol: "tcp", Port: 80}}}}
	return svcs
}
func (APIConnFederationTest) ServiceList() []*object.Service {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	svcs := []*object.Service{{Name: "svc1", Namespace: "testns", ClusterIP: "10.0.0.1", Ports: []api.ServicePort{{Name: "http", Protocol: "tcp", Port: 80}}}, {Name: "hdls1", Namespace: "testns", ClusterIP: api.ClusterIPNone}, {Name: "external", Namespace: "testns", ExternalName: "ext.interwebs.test", Ports: []api.ServicePort{{Name: "http", Protocol: "tcp", Port: 80}}}}
	return svcs
}
func (APIConnFederationTest) EpIndex(string) []*object.Endpoints {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	eps := []*object.Endpoints{{Subsets: []object.EndpointSubset{{Addresses: []object.EndpointAddress{{IP: "172.0.0.1", Hostname: "ep1a"}}, Ports: []object.EndpointPort{{Port: 80, Protocol: "tcp", Name: "http"}}}}, Name: "svc1", Namespace: "testns"}}
	return eps
}
func (APIConnFederationTest) EndpointsList() []*object.Endpoints {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	eps := []*object.Endpoints{{Subsets: []object.EndpointSubset{{Addresses: []object.EndpointAddress{{IP: "172.0.0.1", Hostname: "ep1a"}}, Ports: []object.EndpointPort{{Port: 80, Protocol: "tcp", Name: "http"}}}}, Name: "svc1", Namespace: "testns"}}
	return eps
}
func (a APIConnFederationTest) GetNodeByName(name string) (*api.Node, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &api.Node{ObjectMeta: meta.ObjectMeta{Name: "test.node.foo.bar", Labels: map[string]string{kubernetes.LabelRegion: a.region, kubernetes.LabelZone: a.zone}}}, nil
}
func (APIConnFederationTest) GetNamespaceByName(name string) (*api.Namespace, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &api.Namespace{ObjectMeta: meta.ObjectMeta{Name: name}}, nil
}
