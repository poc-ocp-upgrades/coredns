package kubernetes

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnsutil"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/plugin/pkg/parse"
	"github.com/coredns/coredns/plugin/pkg/upstream"
	"github.com/mholt/caddy"
	"github.com/miekg/dns"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	_ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
	"k8s.io/client-go/tools/clientcmd"
)

var log = clog.NewWithPlugin("kubernetes")

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	flag.Set("logtostderr", "true")
	caddy.RegisterPlugin("kubernetes", caddy.Plugin{ServerType: "dns", Action: setup})
}
func setup(c *caddy.Controller) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	os.Stderr = os.Stdout
	k, err := kubernetesParse(c)
	if err != nil {
		return plugin.Error("kubernetes", err)
	}
	err = k.InitKubeCache()
	if err != nil {
		return plugin.Error("kubernetes", err)
	}
	k.RegisterKubeCache(c)
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		k.Next = next
		return k
	})
	return nil
}
func (k *Kubernetes) RegisterKubeCache(c *caddy.Controller) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.OnStartup(func() error {
		go k.APIConn.Run()
		if k.APIProxy != nil {
			k.APIProxy.Run()
		}
		timeout := time.After(5 * time.Second)
		ticker := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				if k.APIConn.HasSynced() {
					return nil
				}
			case <-timeout:
				return nil
			}
		}
	})
	c.OnShutdown(func() error {
		if k.APIProxy != nil {
			k.APIProxy.Stop()
		}
		return k.APIConn.Stop()
	})
}
func kubernetesParse(c *caddy.Controller) (*Kubernetes, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var (
		k8s	*Kubernetes
		err	error
	)
	i := 0
	for c.Next() {
		if i > 0 {
			return nil, plugin.ErrOnce
		}
		i++
		k8s, err = ParseStanza(c)
		if err != nil {
			return k8s, err
		}
	}
	return k8s, nil
}
func ParseStanza(c *caddy.Controller) (*Kubernetes, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	k8s := New([]string{""})
	k8s.interfaceAddrsFunc = localPodIP
	k8s.autoPathSearch = searchFromResolvConf()
	opts := dnsControlOpts{initEndpointsCache: true, ignoreEmptyService: false, resyncPeriod: defaultResyncPeriod}
	k8s.opts = opts
	zones := c.RemainingArgs()
	if len(zones) != 0 {
		k8s.Zones = zones
		for i := 0; i < len(k8s.Zones); i++ {
			k8s.Zones[i] = plugin.Host(k8s.Zones[i]).Normalize()
		}
	} else {
		k8s.Zones = make([]string, len(c.ServerBlockKeys))
		for i := 0; i < len(c.ServerBlockKeys); i++ {
			k8s.Zones[i] = plugin.Host(c.ServerBlockKeys[i]).Normalize()
		}
	}
	k8s.primaryZoneIndex = -1
	for i, z := range k8s.Zones {
		if dnsutil.IsReverse(z) > 0 {
			continue
		}
		k8s.primaryZoneIndex = i
		break
	}
	if k8s.primaryZoneIndex == -1 {
		return nil, errors.New("non-reverse zone name must be used")
	}
	for c.NextBlock() {
		switch c.Val() {
		case "endpoint_pod_names":
			args := c.RemainingArgs()
			if len(args) > 0 {
				return nil, c.ArgErr()
			}
			k8s.endpointNameMode = true
			continue
		case "pods":
			args := c.RemainingArgs()
			if len(args) == 1 {
				switch args[0] {
				case podModeDisabled, podModeInsecure, podModeVerified:
					k8s.podMode = args[0]
				default:
					return nil, fmt.Errorf("wrong value for pods: %s,  must be one of: disabled, verified, insecure", args[0])
				}
				continue
			}
			return nil, c.ArgErr()
		case "namespaces":
			args := c.RemainingArgs()
			if len(args) > 0 {
				for _, a := range args {
					k8s.Namespaces[a] = struct{}{}
				}
				continue
			}
			return nil, c.ArgErr()
		case "endpoint":
			args := c.RemainingArgs()
			if len(args) > 0 {
				k8s.APIServerList = args
				if len(args) > 1 {
					for i := range args {
						parts := strings.SplitN(args[i], "://", 2)
						if len(parts) == 2 && parts[0] != "http" {
							return nil, fmt.Errorf("multiple endpoints can only accept http, found: %v", args[i])
						}
					}
				}
				continue
			}
			return nil, c.ArgErr()
		case "tls":
			args := c.RemainingArgs()
			if len(args) == 3 {
				k8s.APIClientCert, k8s.APIClientKey, k8s.APICertAuth = args[0], args[1], args[2]
				continue
			}
			return nil, c.ArgErr()
		case "resyncperiod":
			args := c.RemainingArgs()
			if len(args) > 0 {
				rp, err := time.ParseDuration(args[0])
				if err != nil {
					return nil, fmt.Errorf("unable to parse resync duration value: '%v': %v", args[0], err)
				}
				k8s.opts.resyncPeriod = rp
				continue
			}
			return nil, c.ArgErr()
		case "labels":
			args := c.RemainingArgs()
			if len(args) > 0 {
				labelSelectorString := strings.Join(args, " ")
				ls, err := meta.ParseToLabelSelector(labelSelectorString)
				if err != nil {
					return nil, fmt.Errorf("unable to parse label selector value: '%v': %v", labelSelectorString, err)
				}
				k8s.opts.labelSelector = ls
				continue
			}
			return nil, c.ArgErr()
		case "fallthrough":
			k8s.Fall.SetZonesFromArgs(c.RemainingArgs())
		case "upstream":
			args := c.RemainingArgs()
			u, err := upstream.New(args)
			if err != nil {
				return nil, err
			}
			k8s.Upstream = u
		case "ttl":
			args := c.RemainingArgs()
			if len(args) == 0 {
				return nil, c.ArgErr()
			}
			t, err := strconv.Atoi(args[0])
			if err != nil {
				return nil, err
			}
			if t < 0 || t > 3600 {
				return nil, c.Errf("ttl must be in range [0, 3600]: %d", t)
			}
			k8s.ttl = uint32(t)
		case "transfer":
			tos, froms, err := parse.Transfer(c, false)
			if err != nil {
				return nil, err
			}
			if len(froms) != 0 {
				return nil, c.Errf("transfer from is not supported with this plugin")
			}
			k8s.TransferTo = tos
		case "noendpoints":
			if len(c.RemainingArgs()) != 0 {
				return nil, c.ArgErr()
			}
			k8s.opts.initEndpointsCache = false
		case "ignore":
			args := c.RemainingArgs()
			if len(args) > 0 {
				ignore := args[0]
				if ignore == "empty_service" {
					k8s.opts.ignoreEmptyService = true
					continue
				} else {
					return nil, fmt.Errorf("unable to parse ignore value: '%v'", ignore)
				}
			}
		case "kubeconfig":
			args := c.RemainingArgs()
			if len(args) == 2 {
				config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(&clientcmd.ClientConfigLoadingRules{ExplicitPath: args[0]}, &clientcmd.ConfigOverrides{CurrentContext: args[1]})
				k8s.ClientConfig = config
				continue
			}
			return nil, c.ArgErr()
		default:
			return nil, c.Errf("unknown property '%s'", c.Val())
		}
	}
	return k8s, nil
}
func searchFromResolvConf() []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	rc, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return nil
	}
	plugin.Zones(rc.Search).Normalize()
	return rc.Search
}

const defaultResyncPeriod = 5 * time.Minute
