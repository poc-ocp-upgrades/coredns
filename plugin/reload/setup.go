package reload

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/mholt/caddy"
)

var log = clog.NewWithPlugin("reload")

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	caddy.RegisterPlugin("reload", caddy.Plugin{ServerType: "dns", Action: setup})
}

var r = reload{interval: defaultInterval, usage: unused, quit: make(chan bool)}
var once sync.Once
var shutOnce sync.Once

func setup(c *caddy.Controller) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.Next()
	args := c.RemainingArgs()
	if len(args) > 2 {
		return plugin.Error("reload", c.ArgErr())
	}
	i := defaultInterval
	if len(args) > 0 {
		d, err := time.ParseDuration(args[0])
		if err != nil {
			return plugin.Error("reload", err)
		}
		i = d
	}
	if i < minInterval {
		return plugin.Error("reload", fmt.Errorf("interval value must be greater or equal to %v", minInterval))
	}
	j := defaultJitter
	if len(args) > 1 {
		d, err := time.ParseDuration(args[1])
		if err != nil {
			return plugin.Error("reload", err)
		}
		j = d
	}
	if j < minJitter {
		return plugin.Error("reload", fmt.Errorf("jitter value must be greater or equal to %v", minJitter))
	}
	if j > i/2 {
		j = i / 2
	}
	jitter := time.Duration(rand.Int63n(j.Nanoseconds()) - (j.Nanoseconds() / 2))
	i = i + jitter
	r.interval = i
	r.usage = used
	once.Do(func() {
		caddy.RegisterEventHook("reload", hook)
	})
	shutOnce.Do(func() {
		c.OnFinalShutdown(func() error {
			r.quit <- true
			return nil
		})
	})
	return nil
}

const (
	minJitter	= 1 * time.Second
	minInterval	= 2 * time.Second
	defaultInterval	= 30 * time.Second
	defaultJitter	= 15 * time.Second
)
