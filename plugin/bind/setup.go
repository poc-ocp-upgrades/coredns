package bind

import (
	"fmt"
	"net"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/mholt/caddy"
)

func setup(c *caddy.Controller) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	config := dnsserver.GetConfig(c)
	all := []string{}
	for c.Next() {
		addrs := c.RemainingArgs()
		if len(addrs) == 0 {
			return plugin.Error("bind", fmt.Errorf("at least one address is expected"))
		}
		for _, addr := range addrs {
			if net.ParseIP(addr) == nil {
				return plugin.Error("bind", fmt.Errorf("not a valid IP address: %s", addr))
			}
		}
		all = append(all, addrs...)
	}
	config.ListenHosts = all
	return nil
}
