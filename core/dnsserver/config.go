package dnsserver

import (
	"crypto/tls"
	"fmt"
	"github.com/coredns/coredns/plugin"
	"github.com/mholt/caddy"
)

type Config struct {
	Zone		string
	ListenHosts	[]string
	Port		string
	Root		string
	Debug		bool
	Transport	string
	FilterFunc	func(string) bool
	TLSConfig	*tls.Config
	Plugin		[]plugin.Plugin
	pluginChain	plugin.Handler
	registry	map[string]plugin.Handler
}

func keyForConfig(blocIndex int, blocKeyIndex int) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%d:%d", blocIndex, blocKeyIndex)
}
func GetConfig(c *caddy.Controller) *Config {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ctx := c.Context().(*dnsContext)
	key := keyForConfig(c.ServerBlockIndex, c.ServerBlockKeyIndex)
	if cfg, ok := ctx.keysToConfigs[key]; ok {
		return cfg
	}
	ctx.saveConfig(key, &Config{ListenHosts: []string{""}})
	return GetConfig(c)
}
