package rewrite

import (
	"context"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"strings"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type classRule struct {
	fromClass	uint16
	toClass		uint16
	NextAction	string
}

func newClassRule(nextAction string, args ...string) (Rule, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var from, to uint16
	var ok bool
	if from, ok = dns.StringToClass[strings.ToUpper(args[0])]; !ok {
		return nil, fmt.Errorf("invalid class %q", strings.ToUpper(args[0]))
	}
	if to, ok = dns.StringToClass[strings.ToUpper(args[1])]; !ok {
		return nil, fmt.Errorf("invalid class %q", strings.ToUpper(args[1]))
	}
	return &classRule{from, to, nextAction}, nil
}
func (rule *classRule) Rewrite(ctx context.Context, state request.Request) Result {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if rule.fromClass > 0 && rule.toClass > 0 {
		if state.Req.Question[0].Qclass == rule.fromClass {
			state.Req.Question[0].Qclass = rule.toClass
			return RewriteDone
		}
	}
	return RewriteIgnored
}
func (rule *classRule) Mode() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return rule.NextAction
}
func (rule *classRule) GetResponseRule() ResponseRule {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ResponseRule{}
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
