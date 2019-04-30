package dnsserver

import (
	"flag"
	"fmt"
	"net"
	"strings"
	"time"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnsutil"
	"github.com/coredns/coredns/plugin/pkg/parse"
	"github.com/coredns/coredns/plugin/pkg/transport"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyfile"
)

const serverType = "dns"

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	flag.StringVar(&Port, serverType+".port", DefaultPort, "Default port")
	caddy.RegisterServerType(serverType, caddy.ServerType{Directives: func() []string {
		return Directives
	}, DefaultInput: func() caddy.Input {
		return caddy.CaddyfileInput{Filepath: "Corefile", Contents: []byte(".:" + Port + " {\nwhoami\n}\n"), ServerTypeName: serverType}
	}, NewContext: newContext})
}
func newContext(i *caddy.Instance) caddy.Context {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &dnsContext{keysToConfigs: make(map[string]*Config)}
}

type dnsContext struct {
	keysToConfigs	map[string]*Config
	configs		[]*Config
}

func (h *dnsContext) saveConfig(key string, cfg *Config) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	h.configs = append(h.configs, cfg)
	h.keysToConfigs[key] = cfg
}
func (h *dnsContext) InspectServerBlocks(sourceFile string, serverBlocks []caddyfile.ServerBlock) ([]caddyfile.ServerBlock, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for ib, s := range serverBlocks {
		for ik, k := range s.Keys {
			za, err := normalizeZone(k)
			if err != nil {
				return nil, err
			}
			s.Keys[ik] = za.String()
			cfg := &Config{Zone: za.Zone, ListenHosts: []string{""}, Port: za.Port, Transport: za.Transport}
			keyConfig := keyForConfig(ib, ik)
			if za.IPNet == nil {
				h.saveConfig(keyConfig, cfg)
				continue
			}
			ones, bits := za.IPNet.Mask.Size()
			if (bits-ones)%8 != 0 {
				cfg.FilterFunc = func(s string) bool {
					addr := dnsutil.ExtractAddressFromReverse(strings.ToLower(s))
					if addr == "" {
						return true
					}
					return za.IPNet.Contains(net.ParseIP(addr))
				}
			}
			h.saveConfig(keyConfig, cfg)
		}
	}
	return serverBlocks, nil
}
func (h *dnsContext) MakeServers() ([]caddy.Server, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	errValid := h.validateZonesAndListeningAddresses()
	if errValid != nil {
		return nil, errValid
	}
	groups, err := groupConfigsByListenAddr(h.configs)
	if err != nil {
		return nil, err
	}
	var servers []caddy.Server
	for addr, group := range groups {
		switch tr, _ := parse.Transport(addr); tr {
		case transport.DNS:
			s, err := NewServer(addr, group)
			if err != nil {
				return nil, err
			}
			servers = append(servers, s)
		case transport.TLS:
			s, err := NewServerTLS(addr, group)
			if err != nil {
				return nil, err
			}
			servers = append(servers, s)
		case transport.GRPC:
			s, err := NewServergRPC(addr, group)
			if err != nil {
				return nil, err
			}
			servers = append(servers, s)
		case transport.HTTPS:
			s, err := NewServerHTTPS(addr, group)
			if err != nil {
				return nil, err
			}
			servers = append(servers, s)
		}
	}
	return servers, nil
}
func (c *Config) AddPlugin(m plugin.Plugin) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.Plugin = append(c.Plugin, m)
}
func (c *Config) registerHandler(h plugin.Handler) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.registry == nil {
		c.registry = make(map[string]plugin.Handler)
	}
	c.registry[h.Name()] = h
}
func (c *Config) Handler(name string) plugin.Handler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.registry == nil {
		return nil
	}
	if h, ok := c.registry[name]; ok {
		return h
	}
	return nil
}
func (c *Config) Handlers() []plugin.Handler {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.registry == nil {
		return nil
	}
	hs := make([]plugin.Handler, 0, len(c.registry))
	for k := range c.registry {
		hs = append(hs, c.registry[k])
	}
	return hs
}
func (h *dnsContext) validateZonesAndListeningAddresses() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	checker := newOverlapZone()
	for _, conf := range h.configs {
		for _, h := range conf.ListenHosts {
			akey := zoneAddr{Transport: conf.Transport, Zone: conf.Zone, Address: h, Port: conf.Port}
			existZone, overlapZone := checker.registerAndCheck(akey)
			if existZone != nil {
				return fmt.Errorf("cannot serve %s - it is already defined", akey.String())
			}
			if overlapZone != nil {
				return fmt.Errorf("cannot serve %s - zone overlap listener capacity with %v", akey.String(), overlapZone.String())
			}
		}
	}
	return nil
}
func groupConfigsByListenAddr(configs []*Config) (map[string][]*Config, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	groups := make(map[string][]*Config)
	for _, conf := range configs {
		for _, h := range conf.ListenHosts {
			addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(h, conf.Port))
			if err != nil {
				return nil, err
			}
			addrstr := conf.Transport + "://" + addr.String()
			groups[addrstr] = append(groups[addrstr], conf)
		}
	}
	return groups, nil
}

const DefaultPort = transport.Port

var (
	Port		= DefaultPort
	GracefulTimeout	time.Duration
)
var _ caddy.GracefulServer = new(Server)
