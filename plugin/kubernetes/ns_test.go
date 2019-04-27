package kubernetes

import (
	"testing"
	"github.com/coredns/coredns/plugin/kubernetes/object"
	"github.com/coredns/coredns/plugin/pkg/watch"
	api "k8s.io/api/core/v1"
)

type APIConnTest struct{}

func (APIConnTest) HasSynced() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return true
}
func (APIConnTest) Run() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return
}
func (APIConnTest) Stop() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (APIConnTest) PodIndex(string) []*object.Pod {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (APIConnTest) SvcIndex(string) []*object.Service {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (APIConnTest) SvcIndexReverse(string) []*object.Service {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (APIConnTest) EpIndex(string) []*object.Endpoints {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (APIConnTest) EndpointsList() []*object.Endpoints {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (APIConnTest) Modified() int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return 0
}
func (APIConnTest) SetWatchChan(watch.Chan) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (APIConnTest) Watch(string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (APIConnTest) StopWatching(string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (APIConnTest) ServiceList() []*object.Service {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	svcs := []*object.Service{{Name: "dns-service", Namespace: "kube-system", ClusterIP: "10.0.0.111"}}
	return svcs
}
func (APIConnTest) EpIndexReverse(string) []*object.Endpoints {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	eps := []*object.Endpoints{{Subsets: []object.EndpointSubset{{Addresses: []object.EndpointAddress{{IP: "127.0.0.1"}}}}, Name: "dns-service", Namespace: "kube-system"}}
	return eps
}
func (APIConnTest) GetNodeByName(name string) (*api.Node, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &api.Node{}, nil
}
func (APIConnTest) GetNamespaceByName(name string) (*api.Namespace, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &api.Namespace{}, nil
}
func TestNsAddr(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	k := New([]string{"inter.webs.test."})
	k.APIConn = &APIConnTest{}
	cdr := k.nsAddr()
	expected := "10.0.0.111"
	if cdr.A.String() != expected {
		t.Errorf("Expected A to be %q, got %q", expected, cdr.A.String())
	}
	expected = "dns-service.kube-system.svc."
	if cdr.Hdr.Name != expected {
		t.Errorf("Expected Hdr.Name to be %q, got %q", expected, cdr.Hdr.Name)
	}
}
