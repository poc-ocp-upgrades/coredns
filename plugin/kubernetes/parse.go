package kubernetes

import (
	"github.com/coredns/coredns/plugin/pkg/dnsutil"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type recordRequest struct {
	port		string
	protocol	string
	endpoint	string
	service		string
	namespace	string
	podOrSvc	string
}

func parseRequest(state request.Request) (r recordRequest, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	base, _ := dnsutil.TrimZone(state.Name(), state.Zone)
	if base == "" || base == Svc || base == Pod {
		return r, nil
	}
	segs := dns.SplitDomainName(base)
	r.port = "*"
	r.protocol = "*"
	last := len(segs) - 1
	if last < 0 {
		return r, nil
	}
	r.podOrSvc = segs[last]
	if r.podOrSvc != Pod && r.podOrSvc != Svc {
		return r, errInvalidRequest
	}
	last--
	if last < 0 {
		return r, nil
	}
	r.namespace = segs[last]
	last--
	if last < 0 {
		return r, nil
	}
	r.service = segs[last]
	last--
	if last < 0 {
		return r, nil
	}
	switch last {
	case 0:
		r.endpoint = segs[last]
	case 1:
		r.protocol = stripUnderscore(segs[last])
		r.port = stripUnderscore(segs[last-1])
	default:
		return r, errInvalidRequest
	}
	return r, nil
}
func stripUnderscore(s string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s[0] != '_' {
		return s
	}
	return s[1:]
}
func (r recordRequest) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := r.port
	s += "." + r.protocol
	s += "." + r.endpoint
	s += "." + r.service
	s += "." + r.namespace
	s += "." + r.podOrSvc
	return s
}
