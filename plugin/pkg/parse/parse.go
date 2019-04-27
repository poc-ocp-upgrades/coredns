package parse

import (
	"fmt"
	"github.com/coredns/coredns/plugin/pkg/transport"
	"github.com/mholt/caddy"
)

func Transfer(c *caddy.Controller, secondary bool) (tos, froms []string, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !c.NextArg() {
		return nil, nil, c.ArgErr()
	}
	value := c.Val()
	switch value {
	case "to":
		tos = c.RemainingArgs()
		for i := range tos {
			if tos[i] != "*" {
				normalized, err := HostPort(tos[i], transport.Port)
				if err != nil {
					return nil, nil, err
				}
				tos[i] = normalized
			}
		}
	case "from":
		if !secondary {
			return nil, nil, fmt.Errorf("can't use `transfer from` when not being a secondary")
		}
		froms = c.RemainingArgs()
		for i := range froms {
			if froms[i] != "*" {
				normalized, err := HostPort(froms[i], transport.Port)
				if err != nil {
					return nil, nil, err
				}
				froms[i] = normalized
			} else {
				return nil, nil, fmt.Errorf("can't use '*' in transfer from")
			}
		}
	}
	return
}
