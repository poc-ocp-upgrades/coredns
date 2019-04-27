package kubernetes

func (k *Kubernetes) namespace(n string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ns, err := k.APIConn.GetNamespaceByName(n)
	if err != nil {
		return false
	}
	return ns.ObjectMeta.Name == n
}
func (k *Kubernetes) namespaceExposed(namespace string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, ok := k.Namespaces[namespace]
	if len(k.Namespaces) > 0 && !ok {
		return false
	}
	return true
}
