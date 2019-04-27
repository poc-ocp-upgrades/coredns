package metadata

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Metadata struct {
	Zones		[]string
	Providers	[]Provider
	Next		plugin.Handler
}

func (m *Metadata) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "metadata"
}
func (m *Metadata) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ctx = context.WithValue(ctx, key{}, md{})
	state := request.Request{W: w, Req: r}
	if plugin.Zones(m.Zones).Matches(state.Name()) != "" {
		for _, p := range m.Providers {
			ctx = p.Metadata(ctx, state)
		}
	}
	rcode, err := plugin.NextOrFailure(m.Name(), m.Next, ctx, w, r)
	return rcode, err
}
