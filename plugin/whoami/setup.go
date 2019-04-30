package whoami

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/mholt/caddy"
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	caddy.RegisterPlugin("whoami", caddy.Plugin{ServerType: "dns", Action: setup})
}
func setup(c *caddy.Controller) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.Next()
	if c.NextArg() {
		return plugin.Error("whoami", c.ArgErr())
	}
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Whoami{}
	})
	return nil
}
