package file

func (z *Zone) OnShutdown() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if 0 < z.ReloadInterval {
		z.reloadShutdown <- true
	}
	return nil
}
