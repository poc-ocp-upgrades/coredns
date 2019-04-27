package test

func Fuzz(data []byte) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, _, _, err := CoreDNSServerAndPorts(string(data))
	if err != nil {
		return 1
	}
	return 0
}
