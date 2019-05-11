package proxy

import (
	"context"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Exchanger interface {
	Exchange(ctx context.Context, addr string, state request.Request) (*dns.Msg, error)
	Protocol() string
	Transport() string
	OnStartup(*Proxy) error
	OnShutdown(*Proxy) error
}
