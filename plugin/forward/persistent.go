package forward

import (
	"crypto/tls"
	"net"
	"sort"
	"time"
	"github.com/miekg/dns"
)

type persistConn struct {
	c	*dns.Conn
	used	time.Time
}
type Transport struct {
	avgDialTime	int64
	conns		map[string][]*persistConn
	expire		time.Duration
	addr		string
	tlsConfig	*tls.Config
	dial		chan string
	yield		chan *dns.Conn
	ret		chan *dns.Conn
	stop		chan bool
}

func newTransport(addr string) *Transport {
	_logClusterCodePath()
	defer _logClusterCodePath()
	t := &Transport{avgDialTime: int64(maxDialTimeout / 2), conns: make(map[string][]*persistConn), expire: defaultExpire, addr: addr, dial: make(chan string), yield: make(chan *dns.Conn), ret: make(chan *dns.Conn), stop: make(chan bool)}
	return t
}
func (t *Transport) len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := 0
	for _, conns := range t.conns {
		l += len(conns)
	}
	return l
}
func (t *Transport) connManager() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ticker := time.NewTicker(t.expire)
Wait:
	for {
		select {
		case proto := <-t.dial:
			if stack := t.conns[proto]; len(stack) > 0 {
				pc := stack[len(stack)-1]
				if time.Since(pc.used) < t.expire {
					t.conns[proto] = stack[:len(stack)-1]
					t.ret <- pc.c
					continue Wait
				}
				t.conns[proto] = nil
				go closeConns(stack)
			}
			SocketGauge.WithLabelValues(t.addr).Set(float64(t.len()))
			t.ret <- nil
		case conn := <-t.yield:
			SocketGauge.WithLabelValues(t.addr).Set(float64(t.len() + 1))
			if _, ok := conn.Conn.(*net.UDPConn); ok {
				t.conns["udp"] = append(t.conns["udp"], &persistConn{conn, time.Now()})
				continue Wait
			}
			if t.tlsConfig == nil {
				t.conns["tcp"] = append(t.conns["tcp"], &persistConn{conn, time.Now()})
				continue Wait
			}
			t.conns["tcp-tls"] = append(t.conns["tcp-tls"], &persistConn{conn, time.Now()})
		case <-ticker.C:
			t.cleanup(false)
		case <-t.stop:
			t.cleanup(true)
			close(t.ret)
			return
		}
	}
}
func closeConns(conns []*persistConn) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, pc := range conns {
		pc.c.Close()
	}
}
func (t *Transport) cleanup(all bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	staleTime := time.Now().Add(-t.expire)
	for proto, stack := range t.conns {
		if len(stack) == 0 {
			continue
		}
		if all {
			t.conns[proto] = nil
			go closeConns(stack)
			continue
		}
		if stack[0].used.After(staleTime) {
			continue
		}
		good := sort.Search(len(stack), func(i int) bool {
			return stack[i].used.After(staleTime)
		})
		t.conns[proto] = stack[good:]
		go closeConns(stack[:good])
	}
}
func (t *Transport) Yield(c *dns.Conn) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.yield <- c
}
func (t *Transport) Start() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	go t.connManager()
}
func (t *Transport) Stop() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	close(t.stop)
}
func (t *Transport) SetExpire(expire time.Duration) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.expire = expire
}
func (t *Transport) SetTLSConfig(cfg *tls.Config) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.tlsConfig = cfg
}

const (
	defaultExpire	= 10 * time.Second
	minDialTimeout	= 1 * time.Second
	maxDialTimeout	= 30 * time.Second
	readTimeout	= 2 * time.Second
)
