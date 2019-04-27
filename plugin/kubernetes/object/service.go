package object

import (
	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Service struct {
	Version		string
	Name		string
	Namespace	string
	Index		string
	ClusterIP	string
	Type		api.ServiceType
	ExternalName	string
	Ports		[]api.ServicePort
	ExternalIPs	[]string
	*Empty
}

func ServiceKey(name, namespace string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return name + "." + namespace
}
func ToService(obj interface{}) interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	svc, ok := obj.(*api.Service)
	if !ok {
		return nil
	}
	s := &Service{Version: svc.GetResourceVersion(), Name: svc.GetName(), Namespace: svc.GetNamespace(), Index: ServiceKey(svc.GetName(), svc.GetNamespace()), ClusterIP: svc.Spec.ClusterIP, Type: svc.Spec.Type, ExternalName: svc.Spec.ExternalName, ExternalIPs: make([]string, len(svc.Status.LoadBalancer.Ingress)+len(svc.Spec.ExternalIPs))}
	if len(svc.Spec.Ports) == 0 {
		s.Ports = []api.ServicePort{{Port: -1}}
	} else {
		s.Ports = make([]api.ServicePort, len(svc.Spec.Ports))
		copy(s.Ports, svc.Spec.Ports)
	}
	li := copy(s.ExternalIPs, svc.Spec.ExternalIPs)
	for i, lb := range svc.Status.LoadBalancer.Ingress {
		s.ExternalIPs[li+i] = lb.IP
	}
	*svc = api.Service{}
	return s
}

var _ runtime.Object = &Service{}

func (s *Service) DeepCopyObject() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s1 := &Service{Version: s.Version, Name: s.Name, Namespace: s.Namespace, Index: s.Index, ClusterIP: s.ClusterIP, Type: s.Type, ExternalName: s.ExternalName, Ports: make([]api.ServicePort, len(s.Ports)), ExternalIPs: make([]string, len(s.ExternalIPs))}
	copy(s1.Ports, s.Ports)
	copy(s1.ExternalIPs, s.ExternalIPs)
	return s1
}
func (s *Service) GetNamespace() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return s.Namespace
}
func (s *Service) SetNamespace(namespace string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (s *Service) GetName() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return s.Name
}
func (s *Service) SetName(name string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (s *Service) GetResourceVersion() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return s.Version
}
func (s *Service) SetResourceVersion(version string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
