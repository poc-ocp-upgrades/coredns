package auto

func rewriteToExpand(s string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	copy := ""
	for _, c := range s {
		if c == '{' {
			copy += "$"
		}
		copy += string(c)
	}
	return copy
}
