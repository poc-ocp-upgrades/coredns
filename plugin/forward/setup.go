package forward

import (
	"fmt"
	"strconv"
	"time"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	"github.com/coredns/coredns/plugin/pkg/parse"
	pkgtls "github.com/coredns/coredns/plugin/pkg/tls"
	"github.com/coredns/coredns/plugin/pkg/transport"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyfile"
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	caddy.RegisterPlugin("forward", caddy.Plugin{ServerType: "dns", Action: setup})
}
func setup(c *caddy.Controller) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	f, err := parseForward(c)
	if err != nil {
		return plugin.Error("forward", err)
	}
	if f.Len() > max {
		return plugin.Error("forward", fmt.Errorf("more than %d TOs configured: %d", max, f.Len()))
	}
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		f.Next = next
		return f
	})
	c.OnStartup(func() error {
		metrics.MustRegister(c, RequestCount, RcodeCount, RequestDuration, HealthcheckFailureCount, SocketGauge)
		return f.OnStartup()
	})
	c.OnShutdown(func() error {
		return f.OnShutdown()
	})
	return nil
}
func (f *Forward) OnStartup() (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, p := range f.proxies {
		p.start(f.hcInterval)
	}
	return nil
}
func (f *Forward) OnShutdown() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, p := range f.proxies {
		p.close()
	}
	return nil
}
func (f *Forward) Close() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	f.OnShutdown()
}
func parseForward(c *caddy.Controller) (*Forward, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var (
		f	*Forward
		err	error
		i	int
	)
	for c.Next() {
		if i > 0 {
			return nil, plugin.ErrOnce
		}
		i++
		f, err = ParseForwardStanza(&c.Dispenser)
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}
func ParseForwardStanza(c *caddyfile.Dispenser) (*Forward, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	f := New()
	if !c.Args(&f.from) {
		return f, c.ArgErr()
	}
	f.from = plugin.Host(f.from).Normalize()
	to := c.RemainingArgs()
	if len(to) == 0 {
		return f, c.ArgErr()
	}
	toHosts, err := parse.HostPortOrFile(to...)
	if err != nil {
		return f, err
	}
	transports := make([]string, len(toHosts))
	for i, host := range toHosts {
		trans, h := parse.Transport(host)
		p := NewProxy(h, trans)
		f.proxies = append(f.proxies, p)
		transports[i] = trans
	}
	for c.NextBlock() {
		if err := parseBlock(c, f); err != nil {
			return f, err
		}
	}
	if f.tlsServerName != "" {
		f.tlsConfig.ServerName = f.tlsServerName
	}
	for i := range f.proxies {
		if transports[i] == transport.TLS {
			f.proxies[i].SetTLSConfig(f.tlsConfig)
		}
		f.proxies[i].SetExpire(f.expire)
	}
	return f, nil
}
func parseBlock(c *caddyfile.Dispenser, f *Forward) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch c.Val() {
	case "except":
		ignore := c.RemainingArgs()
		if len(ignore) == 0 {
			return c.ArgErr()
		}
		for i := 0; i < len(ignore); i++ {
			ignore[i] = plugin.Host(ignore[i]).Normalize()
		}
		f.ignored = ignore
	case "max_fails":
		if !c.NextArg() {
			return c.ArgErr()
		}
		n, err := strconv.Atoi(c.Val())
		if err != nil {
			return err
		}
		if n < 0 {
			return fmt.Errorf("max_fails can't be negative: %d", n)
		}
		f.maxfails = uint32(n)
	case "health_check":
		if !c.NextArg() {
			return c.ArgErr()
		}
		dur, err := time.ParseDuration(c.Val())
		if err != nil {
			return err
		}
		if dur < 0 {
			return fmt.Errorf("health_check can't be negative: %d", dur)
		}
		f.hcInterval = dur
	case "force_tcp":
		if c.NextArg() {
			return c.ArgErr()
		}
		f.opts.forceTCP = true
	case "prefer_udp":
		if c.NextArg() {
			return c.ArgErr()
		}
		f.opts.preferUDP = true
	case "tls":
		args := c.RemainingArgs()
		if len(args) > 3 {
			return c.ArgErr()
		}
		tlsConfig, err := pkgtls.NewTLSConfigFromArgs(args...)
		if err != nil {
			return err
		}
		f.tlsConfig = tlsConfig
	case "tls_servername":
		if !c.NextArg() {
			return c.ArgErr()
		}
		f.tlsServerName = c.Val()
	case "expire":
		if !c.NextArg() {
			return c.ArgErr()
		}
		dur, err := time.ParseDuration(c.Val())
		if err != nil {
			return err
		}
		if dur < 0 {
			return fmt.Errorf("expire can't be negative: %s", dur)
		}
		f.expire = dur
	case "policy":
		if !c.NextArg() {
			return c.ArgErr()
		}
		switch x := c.Val(); x {
		case "random":
			f.p = &random{}
		case "round_robin":
			f.p = &roundRobin{}
		case "sequential":
			f.p = &sequential{}
		default:
			return c.Errf("unknown policy '%s'", x)
		}
	default:
		return c.Errf("unknown property '%s'", c.Val())
	}
	return nil
}

const max = 15
