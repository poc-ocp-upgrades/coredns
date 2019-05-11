package request

import (
	"context"
	"net"
	"strings"
	"github.com/coredns/coredns/plugin/pkg/edns"
	"github.com/miekg/dns"
)

type Request struct {
	Req			*dns.Msg
	W			dns.ResponseWriter
	Zone		string
	Context		context.Context
	size		int
	do			*bool
	name		string
	ip			string
	port		string
	family		int
	localPort	string
	localIP		string
}

func (r *Request) NewWithQuestion(name string, typ uint16) Request {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req1 := Request{W: r.W, Req: r.Req.Copy()}
	req1.Req.Question[0] = dns.Question{Name: dns.Fqdn(name), Qclass: dns.ClassINET, Qtype: typ}
	return req1
}
func (r *Request) IP() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.ip != "" {
		return r.ip
	}
	ip, _, err := net.SplitHostPort(r.W.RemoteAddr().String())
	if err != nil {
		r.ip = r.W.RemoteAddr().String()
		return r.ip
	}
	r.ip = ip
	return r.ip
}
func (r *Request) LocalIP() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.localIP != "" {
		return r.localIP
	}
	ip, _, err := net.SplitHostPort(r.W.LocalAddr().String())
	if err != nil {
		r.localIP = r.W.LocalAddr().String()
		return r.localIP
	}
	r.localIP = ip
	return r.localIP
}
func (r *Request) Port() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.port != "" {
		return r.port
	}
	_, port, err := net.SplitHostPort(r.W.RemoteAddr().String())
	if err != nil {
		r.port = "0"
		return r.port
	}
	r.port = port
	return r.port
}
func (r *Request) LocalPort() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.localPort != "" {
		return r.localPort
	}
	_, port, err := net.SplitHostPort(r.W.LocalAddr().String())
	if err != nil {
		r.localPort = "0"
		return r.localPort
	}
	r.localPort = port
	return r.localPort
}
func (r *Request) RemoteAddr() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.W.RemoteAddr().String()
}
func (r *Request) LocalAddr() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.W.LocalAddr().String()
}
func (r *Request) Proto() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return Proto(r.W)
}
func Proto(w dns.ResponseWriter) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if _, ok := w.RemoteAddr().(*net.UDPAddr); ok {
		return "udp"
	}
	if _, ok := w.RemoteAddr().(*net.TCPAddr); ok {
		return "tcp"
	}
	return "udp"
}
func (r *Request) Family() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.family != 0 {
		return r.family
	}
	var a net.IP
	ip := r.W.RemoteAddr()
	if i, ok := ip.(*net.UDPAddr); ok {
		a = i.IP
	}
	if i, ok := ip.(*net.TCPAddr); ok {
		a = i.IP
	}
	if a.To4() != nil {
		r.family = 1
		return r.family
	}
	r.family = 2
	return r.family
}
func (r *Request) Do() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.do != nil {
		return *r.do
	}
	r.do = new(bool)
	if o := r.Req.IsEdns0(); o != nil {
		*r.do = o.Do()
		return *r.do
	}
	*r.do = false
	return false
}
func (r *Request) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return r.Req.Len()
}
func (r *Request) Size() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.size != 0 {
		return r.size
	}
	size := 0
	if o := r.Req.IsEdns0(); o != nil {
		if r.do == nil {
			r.do = new(bool)
		}
		*r.do = o.Do()
		size = int(o.UDPSize())
	}
	size = edns.Size(r.Proto(), size)
	r.size = size
	return size
}
func (r *Request) SizeAndDo(m *dns.Msg) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	o := r.Req.IsEdns0()
	if o == nil {
		return false
	}
	if mo := m.IsEdns0(); mo != nil {
		mo.Hdr.Name = "."
		mo.Hdr.Rrtype = dns.TypeOPT
		mo.SetVersion(0)
		mo.SetUDPSize(o.UDPSize())
		mo.Hdr.Ttl &= 0xff00
		if o.Do() {
			mo.SetDo()
		}
		return true
	}
	o.Hdr.Name = "."
	o.Hdr.Rrtype = dns.TypeOPT
	o.SetVersion(0)
	o.Hdr.Ttl &= 0xff00
	if len(o.Option) > 0 {
		o.Option = supportedOptions(o.Option)
	}
	m.Extra = append(m.Extra, o)
	return true
}
func (r *Request) Scrub(reply *dns.Msg) *dns.Msg {
	_logClusterCodePath()
	defer _logClusterCodePath()
	size := r.Size()
	reply.Compress = false
	rl := reply.Len()
	if size >= rl {
		if r.Proto() != "udp" {
			return reply
		}
		if rl > 1480 && r.Family() == 1 {
			reply.Compress = true
		}
		if rl > 1220 && r.Family() == 2 {
			reply.Compress = true
		}
		return reply
	}
	reply.Compress = true
	rl = reply.Len()
	if size >= rl {
		return reply
	}
	re := len(reply.Extra)
	if r.Req.IsEdns0() != nil {
		size -= optLen
		re--
	}
	l, m := 0, 0
	origExtra := reply.Extra
	for l <= re {
		m = (l + re) / 2
		reply.Extra = origExtra[:m]
		rl = reply.Len()
		if rl < size {
			l = m + 1
			continue
		}
		if rl > size {
			re = m - 1
			continue
		}
		if rl == size {
			break
		}
	}
	if rl > size && m > 0 {
		reply.Extra = origExtra[:m-1]
		rl = reply.Len()
	}
	if rl <= size {
		return reply
	}
	ra := len(reply.Answer)
	l, m = 0, 0
	origAnswer := reply.Answer
	for l <= ra {
		m = (l + ra) / 2
		reply.Answer = origAnswer[:m]
		rl = reply.Len()
		if rl < size {
			l = m + 1
			continue
		}
		if rl > size {
			ra = m - 1
			continue
		}
		if rl == size {
			break
		}
	}
	if rl > size && m > 0 {
		reply.Answer = origAnswer[:m-1]
	}
	reply.Truncated = true
	return reply
}
func (r *Request) Type() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.Req == nil {
		return ""
	}
	if len(r.Req.Question) == 0 {
		return ""
	}
	return dns.Type(r.Req.Question[0].Qtype).String()
}
func (r *Request) QType() uint16 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.Req == nil {
		return 0
	}
	if len(r.Req.Question) == 0 {
		return 0
	}
	return r.Req.Question[0].Qtype
}
func (r *Request) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.name != "" {
		return r.name
	}
	if r.Req == nil {
		r.name = "."
		return "."
	}
	if len(r.Req.Question) == 0 {
		r.name = "."
		return "."
	}
	r.name = strings.ToLower(dns.Name(r.Req.Question[0].Name).String())
	return r.name
}
func (r *Request) QName() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.Req == nil {
		return "."
	}
	if len(r.Req.Question) == 0 {
		return "."
	}
	return dns.Name(r.Req.Question[0].Name).String()
}
func (r *Request) Class() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.Req == nil {
		return ""
	}
	if len(r.Req.Question) == 0 {
		return ""
	}
	return dns.Class(r.Req.Question[0].Qclass).String()
}
func (r *Request) QClass() uint16 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.Req == nil {
		return 0
	}
	if len(r.Req.Question) == 0 {
		return 0
	}
	return r.Req.Question[0].Qclass
}
func (r *Request) ErrorMessage(rcode int) *dns.Msg {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := new(dns.Msg)
	m.SetRcode(r.Req, rcode)
	return m
}
func (r *Request) Clear() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.name = ""
	r.ip = ""
	r.localIP = ""
	r.port = ""
	r.localPort = ""
	r.family = 0
}
func (r *Request) Match(reply *dns.Msg) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(reply.Question) != 1 {
		return false
	}
	if reply.Response == false {
		return false
	}
	if strings.ToLower(reply.Question[0].Name) != r.Name() {
		return false
	}
	if reply.Question[0].Qtype != r.QType() {
		return false
	}
	return true
}

const optLen = 12
