package dnstest

import (
	"time"
	"github.com/miekg/dns"
)

type Recorder struct {
	dns.ResponseWriter
	Rcode	int
	Len	int
	Msg	*dns.Msg
	Start	time.Time
}

func NewRecorder(w dns.ResponseWriter) *Recorder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &Recorder{ResponseWriter: w, Rcode: 0, Msg: nil, Start: time.Now()}
}
func (r *Recorder) WriteMsg(res *dns.Msg) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.Rcode = res.Rcode
	r.Len += res.Len()
	r.Msg = res
	return r.ResponseWriter.WriteMsg(res)
}
func (r *Recorder) Write(buf []byte) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	n, err := r.ResponseWriter.Write(buf)
	if err == nil {
		r.Len += n
	}
	return n, err
}
