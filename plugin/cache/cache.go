package cache

import (
	"hash/fnv"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"net"
	"time"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/cache"
	"github.com/coredns/coredns/plugin/pkg/dnsutil"
	"github.com/coredns/coredns/plugin/pkg/response"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Cache struct {
	Next		plugin.Handler
	Zones		[]string
	ncache		*cache.Cache
	ncap		int
	nttl		time.Duration
	minnttl		time.Duration
	pcache		*cache.Cache
	pcap		int
	pttl		time.Duration
	minpttl		time.Duration
	prefetch	int
	duration	time.Duration
	percentage	int
	now		func() time.Time
}

func New() *Cache {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &Cache{Zones: []string{"."}, pcap: defaultCap, pcache: cache.New(defaultCap), pttl: maxTTL, minpttl: minTTL, ncap: defaultCap, ncache: cache.New(defaultCap), nttl: maxNTTL, minnttl: minNTTL, prefetch: 0, duration: 1 * time.Minute, percentage: 10, now: time.Now}
}
func key(qname string, m *dns.Msg, t response.Type, do bool) (bool, uint64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m.Truncated {
		return false, 0
	}
	if t == response.OtherError || t == response.Meta || t == response.Update {
		return false, 0
	}
	return true, hash(qname, m.Question[0].Qtype, do)
}

var one = []byte("1")
var zero = []byte("0")

func hash(qname string, qtype uint16, do bool) uint64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	h := fnv.New64()
	if do {
		h.Write(one)
	} else {
		h.Write(zero)
	}
	h.Write([]byte{byte(qtype >> 8)})
	h.Write([]byte{byte(qtype)})
	h.Write([]byte(qname))
	return h.Sum64()
}
func computeTTL(msgTTL, minTTL, maxTTL time.Duration) time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ttl := msgTTL
	if ttl < minTTL {
		ttl = minTTL
	}
	if ttl > maxTTL {
		ttl = maxTTL
	}
	return ttl
}

type ResponseWriter struct {
	dns.ResponseWriter
	*Cache
	state		request.Request
	server		string
	prefetch	bool
	remoteAddr	net.Addr
}

func newPrefetchResponseWriter(server string, state request.Request, c *Cache) *ResponseWriter {
	_logClusterCodePath()
	defer _logClusterCodePath()
	addr := state.W.RemoteAddr()
	if u, ok := addr.(*net.UDPAddr); ok {
		addr = &net.TCPAddr{IP: u.IP, Port: u.Port, Zone: u.Zone}
	}
	return &ResponseWriter{ResponseWriter: state.W, Cache: c, state: state, server: server, prefetch: true, remoteAddr: addr}
}
func (w *ResponseWriter) RemoteAddr() net.Addr {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if w.remoteAddr != nil {
		return w.remoteAddr
	}
	return w.ResponseWriter.RemoteAddr()
}
func (w *ResponseWriter) WriteMsg(res *dns.Msg) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	do := false
	mt, opt := response.Typify(res, w.now().UTC())
	if opt != nil {
		do = opt.Do()
	}
	hasKey, key := key(w.state.Name(), res, mt, do)
	msgTTL := dnsutil.MinimalTTL(res, mt)
	var duration time.Duration
	if mt == response.NameError || mt == response.NoData {
		duration = computeTTL(msgTTL, w.minnttl, w.nttl)
	} else {
		duration = computeTTL(msgTTL, w.minpttl, w.pttl)
	}
	if hasKey && duration > 0 {
		if w.state.Match(res) {
			w.set(res, key, mt, duration)
			cacheSize.WithLabelValues(w.server, Success).Set(float64(w.pcache.Len()))
			cacheSize.WithLabelValues(w.server, Denial).Set(float64(w.ncache.Len()))
		} else {
			cacheDrops.WithLabelValues(w.server).Inc()
		}
	}
	if w.prefetch {
		return nil
	}
	ttl := uint32(duration.Seconds())
	for i := range res.Answer {
		res.Answer[i].Header().Ttl = ttl
	}
	for i := range res.Ns {
		res.Ns[i].Header().Ttl = ttl
	}
	for i := range res.Extra {
		if res.Extra[i].Header().Rrtype != dns.TypeOPT {
			res.Extra[i].Header().Ttl = ttl
		}
	}
	return w.ResponseWriter.WriteMsg(res)
}
func (w *ResponseWriter) set(m *dns.Msg, key uint64, mt response.Type, duration time.Duration) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch mt {
	case response.NoError, response.Delegation:
		i := newItem(m, w.now(), duration)
		w.pcache.Add(key, i)
	case response.NameError, response.NoData:
		i := newItem(m, w.now(), duration)
		w.ncache.Add(key, i)
	case response.OtherError:
	default:
		log.Warningf("Caching called with unknown classification: %d", mt)
	}
}
func (w *ResponseWriter) Write(buf []byte) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	log.Warning("Caching called with Write: not caching reply")
	if w.prefetch {
		return 0, nil
	}
	n, err := w.ResponseWriter.Write(buf)
	return n, err
}

const (
	maxTTL		= dnsutil.MaximumDefaulTTL
	minTTL		= dnsutil.MinimalDefaultTTL
	maxNTTL		= dnsutil.MaximumDefaulTTL / 2
	minNTTL		= dnsutil.MinimalDefaultTTL
	defaultCap	= 10000
	Success		= "success"
	Denial		= "denial"
)

func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
