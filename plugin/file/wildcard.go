package file

import "github.com/miekg/dns"

func replaceWithAsteriskLabel(qname string) (wildcard string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	i, shot := dns.NextLabel(qname, 0)
	if shot {
		return ""
	}
	return "*." + qname[i:]
}
