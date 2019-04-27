package rewrite

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/mholt/caddy"
)

var log = clog.NewWithPlugin("rewrite")

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	caddy.RegisterPlugin("rewrite", caddy.Plugin{ServerType: "dns", Action: setup})
}
func setup(c *caddy.Controller) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	rewrites, err := rewriteParse(c)
	if err != nil {
		return plugin.Error("rewrite", err)
	}
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Rewrite{Next: next, Rules: rewrites}
	})
	return nil
}
func rewriteParse(c *caddy.Controller) ([]Rule, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var rules []Rule
	for c.Next() {
		args := c.RemainingArgs()
		if len(args) < 2 {
			for c.NextBlock() {
				args = append(args, c.Val())
			}
		}
		rule, err := newRule(args...)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}
