package object

import (
	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Pod struct {
	Version		string
	PodIP		string
	Name		string
	Namespace	string
	Deleting	bool
	*Empty
}

func ToPod(obj interface{}) interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pod, ok := obj.(*api.Pod)
	if !ok {
		return nil
	}
	p := &Pod{Version: pod.GetResourceVersion(), PodIP: pod.Status.PodIP, Namespace: pod.GetNamespace(), Name: pod.GetName()}
	t := pod.ObjectMeta.DeletionTimestamp
	if t != nil {
		p.Deleting = !(*t).Time.IsZero()
	}
	*pod = api.Pod{}
	return p
}

var _ runtime.Object = &Pod{}

func (p *Pod) DeepCopyObject() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	p1 := &Pod{Version: p.Version, PodIP: p.PodIP, Namespace: p.Namespace, Name: p.Name, Deleting: p.Deleting}
	return p1
}
func (p *Pod) GetNamespace() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return p.Namespace
}
func (p *Pod) SetNamespace(namespace string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (p *Pod) GetName() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return p.Name
}
func (p *Pod) SetName(name string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (p *Pod) GetResourceVersion() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return p.Version
}
func (p *Pod) SetResourceVersion(version string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
