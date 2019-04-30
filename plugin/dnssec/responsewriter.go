package dnssec

import (
	"time"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type ResponseWriter struct {
	dns.ResponseWriter
	d	Dnssec
	server	string
}

func (d *ResponseWriter) WriteMsg(res *dns.Msg) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	state := request.Request{W: d.ResponseWriter, Req: res}
	zone := plugin.Zones(d.d.zones).Matches(state.Name())
	if zone == "" {
		return d.ResponseWriter.WriteMsg(res)
	}
	state.Zone = zone
	res = d.d.Sign(state, time.Now().UTC(), d.server)
	cacheSize.WithLabelValues(d.server, "signature").Set(float64(d.d.cache.Len()))
	return d.ResponseWriter.WriteMsg(res)
}
func (d *ResponseWriter) Write(buf []byte) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	log.Warning("Dnssec called with Write: not signing reply")
	n, err := d.ResponseWriter.Write(buf)
	return n, err
}
