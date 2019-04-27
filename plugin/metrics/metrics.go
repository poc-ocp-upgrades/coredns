package metrics

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics/vars"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	Next		plugin.Handler
	Addr		string
	Reg		*prometheus.Registry
	ln		net.Listener
	lnSetup		bool
	mux		*http.ServeMux
	srv		*http.Server
	zoneNames	[]string
	zoneMap		map[string]struct{}
	zoneMu		sync.RWMutex
}

func New(addr string) *Metrics {
	_logClusterCodePath()
	defer _logClusterCodePath()
	met := &Metrics{Addr: addr, Reg: prometheus.NewRegistry(), zoneMap: make(map[string]struct{})}
	met.MustRegister(prometheus.NewGoCollector())
	met.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	met.MustRegister(buildInfo)
	met.MustRegister(vars.Panic)
	met.MustRegister(vars.RequestCount)
	met.MustRegister(vars.RequestDuration)
	met.MustRegister(vars.RequestSize)
	met.MustRegister(vars.RequestDo)
	met.MustRegister(vars.RequestType)
	met.MustRegister(vars.ResponseSize)
	met.MustRegister(vars.ResponseRcode)
	return met
}
func (m *Metrics) MustRegister(c prometheus.Collector) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err := m.Reg.Register(c)
	if err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			log.Fatalf("Cannot register metrics collector: %s", err)
		}
	}
}
func (m *Metrics) AddZone(z string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.zoneMu.Lock()
	m.zoneMap[z] = struct{}{}
	m.zoneNames = keys(m.zoneMap)
	m.zoneMu.Unlock()
}
func (m *Metrics) RemoveZone(z string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.zoneMu.Lock()
	delete(m.zoneMap, z)
	m.zoneNames = keys(m.zoneMap)
	m.zoneMu.Unlock()
}
func (m *Metrics) ZoneNames() []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.zoneMu.RLock()
	s := m.zoneNames
	m.zoneMu.RUnlock()
	return s
}
func (m *Metrics) OnStartup() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ln, err := net.Listen("tcp", m.Addr)
	if err != nil {
		log.Errorf("Failed to start metrics handler: %s", err)
		return err
	}
	m.ln = ln
	m.lnSetup = true
	ListenAddr = m.ln.Addr().String()
	m.mux = http.NewServeMux()
	m.mux.Handle("/metrics", promhttp.HandlerFor(m.Reg, promhttp.HandlerOpts{}))
	m.srv = &http.Server{Handler: m.mux}
	go func() {
		m.srv.Serve(m.ln)
	}()
	return nil
}
func (m *Metrics) OnRestart() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !m.lnSetup {
		return nil
	}
	uniqAddr.Unset(m.Addr)
	return m.stopServer()
}
func (m *Metrics) stopServer() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !m.lnSetup {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := m.srv.Shutdown(ctx); err != nil {
		log.Infof("Failed to stop prometheus http server: %s", err)
		return err
	}
	m.lnSetup = false
	m.ln.Close()
	return nil
}
func (m *Metrics) OnFinalShutdown() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.stopServer()
}
func keys(m map[string]struct{}) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sx := []string{}
	for k := range m {
		sx = append(sx, k)
	}
	return sx
}

var ListenAddr string

const shutdownTimeout time.Duration = time.Second * 5

var buildInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: plugin.Namespace, Name: "build_info", Help: "A metric with a constant '1' value labeled by version, revision, and goversion from which CoreDNS was built."}, []string{"version", "revision", "goversion"})
