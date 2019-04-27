package dnstap

import (
	"strings"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/dnstap/dnstapio"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/plugin/pkg/parse"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyfile"
)

var log = clog.NewWithPlugin("dnstap")

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	caddy.RegisterPlugin("dnstap", caddy.Plugin{ServerType: "dns", Action: wrapSetup})
}
func wrapSetup(c *caddy.Controller) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := setup(c); err != nil {
		return plugin.Error("dnstap", err)
	}
	return nil
}

type config struct {
	target	string
	socket	bool
	full	bool
}

func parseConfig(d *caddyfile.Dispenser) (c config, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	d.Next()
	if !d.Args(&c.target) {
		return c, d.ArgErr()
	}
	if strings.HasPrefix(c.target, "tcp://") {
		servers, err := parse.HostPortOrFile(c.target[6:])
		if err != nil {
			return c, d.ArgErr()
		}
		c.target = servers[0]
	} else {
		if strings.HasPrefix(c.target, "unix://") {
			c.target = c.target[7:]
		}
		c.socket = true
	}
	c.full = d.NextArg() && d.Val() == "full"
	return
}
func setup(c *caddy.Controller) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	conf, err := parseConfig(&c.Dispenser)
	if err != nil {
		return err
	}
	dio := dnstapio.New(conf.target, conf.socket)
	dnstap := Dnstap{IO: dio, JoinRawMessage: conf.full}
	c.OnStartup(func() error {
		dio.Connect()
		return nil
	})
	c.OnRestart(func() error {
		dio.Close()
		return nil
	})
	c.OnFinalShutdown(func() error {
		dio.Close()
		return nil
	})
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		dnstap.Next = next
		return dnstap
	})
	return nil
}
