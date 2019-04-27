package request

import "github.com/miekg/dns"

type ScrubWriter struct {
	dns.ResponseWriter
	req	*dns.Msg
}

func NewScrubWriter(req *dns.Msg, w dns.ResponseWriter) *ScrubWriter {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &ScrubWriter{w, req}
}
func (s *ScrubWriter) WriteMsg(m *dns.Msg) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	state := Request{Req: s.req, W: s.ResponseWriter}
	n := state.Scrub(m)
	state.SizeAndDo(n)
	return s.ResponseWriter.WriteMsg(n)
}
