package kubernetes

func (k *Kubernetes) Health() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return k.APIConn.HasSynced()
}
