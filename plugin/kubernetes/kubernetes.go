package kubernetes

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync/atomic"
	"time"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/plugin/kubernetes/object"
	"github.com/coredns/coredns/plugin/pkg/dnsutil"
	"github.com/coredns/coredns/plugin/pkg/fall"
	"github.com/coredns/coredns/plugin/pkg/healthcheck"
	"github.com/coredns/coredns/plugin/pkg/upstream"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	api "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type Kubernetes struct {
	Next			plugin.Handler
	Zones			[]string
	Upstream		upstream.Upstream
	APIServerList		[]string
	APIProxy		*apiProxy
	APICertAuth		string
	APIClientCert		string
	APIClientKey		string
	ClientConfig		clientcmd.ClientConfig
	APIConn			dnsController
	Namespaces		map[string]struct{}
	podMode			string
	endpointNameMode	bool
	Fall			fall.F
	ttl			uint32
	opts			dnsControlOpts
	primaryZoneIndex	int
	interfaceAddrsFunc	func() net.IP
	autoPathSearch		[]string
	TransferTo		[]string
}

func New(zones []string) *Kubernetes {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	k := new(Kubernetes)
	k.Zones = zones
	k.Namespaces = make(map[string]struct{})
	k.interfaceAddrsFunc = func() net.IP {
		return net.ParseIP("127.0.0.1")
	}
	k.podMode = podModeDisabled
	k.ttl = defaultTTL
	return k
}

const (
	podModeDisabled		= "disabled"
	podModeVerified		= "verified"
	podModeInsecure		= "insecure"
	DNSSchemaVersion	= "1.0.1"
	Svc			= "svc"
	Pod			= "pod"
	defaultTTL		= 5
)

var (
	errNoItems		= errors.New("no items found")
	errNsNotExposed		= errors.New("namespace is not exposed")
	errInvalidRequest	= errors.New("invalid query name")
)

func (k *Kubernetes) Services(state request.Request, exact bool, opt plugin.Options) (svcs []msg.Service, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch state.QType() {
	case dns.TypeTXT:
		t, _ := dnsutil.TrimZone(state.Name(), state.Zone)
		segs := dns.SplitDomainName(t)
		if len(segs) != 1 {
			return nil, nil
		}
		if segs[0] != "dns-version" {
			return nil, nil
		}
		svc := msg.Service{Text: DNSSchemaVersion, TTL: 28800, Key: msg.Path(state.QName(), coredns)}
		return []msg.Service{svc}, nil
	case dns.TypeNS:
		ns := k.nsAddr()
		svc := msg.Service{Host: ns.A.String(), Key: msg.Path(state.QName(), coredns)}
		return []msg.Service{svc}, nil
	}
	if state.QType() == dns.TypeA && isDefaultNS(state.Name(), state.Zone) {
		ns := k.nsAddr()
		svc := msg.Service{Host: ns.A.String(), Key: msg.Path(state.QName(), coredns)}
		return []msg.Service{svc}, nil
	}
	s, e := k.Records(state, false)
	if state.QType() != dns.TypeSRV {
		return s, e
	}
	internal := []msg.Service{}
	for _, svc := range s {
		if t, _ := svc.HostType(); t != dns.TypeCNAME {
			internal = append(internal, svc)
		}
	}
	return internal, e
}
func (k *Kubernetes) primaryZone() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return k.Zones[k.primaryZoneIndex]
}
func (k *Kubernetes) Lookup(state request.Request, name string, typ uint16) (*dns.Msg, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return k.Upstream.Lookup(state, name, typ)
}
func (k *Kubernetes) IsNameError(err error) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return err == errNoItems || err == errNsNotExposed || err == errInvalidRequest
}
func (k *Kubernetes) getClientConfig() (*rest.Config, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if k.ClientConfig != nil {
		return k.ClientConfig.ClientConfig()
	}
	loadingRules := &clientcmd.ClientConfigLoadingRules{}
	overrides := &clientcmd.ConfigOverrides{}
	clusterinfo := clientcmdapi.Cluster{}
	authinfo := clientcmdapi.AuthInfo{}
	if len(k.APIServerList) == 0 {
		cc, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		cc.ContentType = "application/vnd.kubernetes.protobuf"
		return cc, err
	}
	endpoint := k.APIServerList[0]
	if len(k.APIServerList) > 1 {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return nil, fmt.Errorf("failed to create kubernetes api proxy: %v", err)
		}
		k.APIProxy = &apiProxy{listener: listener, handler: proxyHandler{HealthCheck: healthcheck.HealthCheck{FailTimeout: 3 * time.Second, MaxFails: 1, Path: "/", Interval: 5 * time.Second}}}
		k.APIProxy.handler.Hosts = make([]*healthcheck.UpstreamHost, len(k.APIServerList))
		for i, entry := range k.APIServerList {
			uh := &healthcheck.UpstreamHost{Name: strings.TrimPrefix(entry, "http://"), CheckDown: func(upstream *proxyHandler) healthcheck.UpstreamHostDownFunc {
				return func(uh *healthcheck.UpstreamHost) bool {
					fails := atomic.LoadInt32(&uh.Fails)
					if fails >= upstream.MaxFails && upstream.MaxFails != 0 {
						return true
					}
					return false
				}
			}(&k.APIProxy.handler)}
			k.APIProxy.handler.Hosts[i] = uh
		}
		k.APIProxy.Handler = &k.APIProxy.handler
		endpoint = fmt.Sprintf("http://%s", listener.Addr())
	}
	clusterinfo.Server = endpoint
	if len(k.APICertAuth) > 0 {
		clusterinfo.CertificateAuthority = k.APICertAuth
	}
	if len(k.APIClientCert) > 0 {
		authinfo.ClientCertificate = k.APIClientCert
	}
	if len(k.APIClientKey) > 0 {
		authinfo.ClientKey = k.APIClientKey
	}
	overrides.ClusterInfo = clusterinfo
	overrides.AuthInfo = authinfo
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides)
	cc, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	cc.ContentType = "application/vnd.kubernetes.protobuf"
	return cc, err
}
func (k *Kubernetes) InitKubeCache() (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	config, err := k.getClientConfig()
	if err != nil {
		return err
	}
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes notification controller: %q", err)
	}
	if k.opts.labelSelector != nil {
		var selector labels.Selector
		selector, err = meta.LabelSelectorAsSelector(k.opts.labelSelector)
		if err != nil {
			return fmt.Errorf("unable to create Selector for LabelSelector '%s': %q", k.opts.labelSelector, err)
		}
		k.opts.selector = selector
	}
	k.opts.initPodCache = k.podMode == podModeVerified
	k.opts.zones = k.Zones
	k.opts.endpointNameMode = k.endpointNameMode
	k.APIConn = newdnsController(kubeClient, k.opts)
	return err
}
func (k *Kubernetes) Records(state request.Request, exact bool) ([]msg.Service, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	r, e := parseRequest(state)
	if e != nil {
		return nil, e
	}
	if r.podOrSvc == "" {
		return nil, nil
	}
	if dnsutil.IsReverse(state.Name()) > 0 {
		return nil, errNoItems
	}
	if !wildcard(r.namespace) && !k.namespaceExposed(r.namespace) {
		return nil, errNsNotExposed
	}
	if r.podOrSvc == Pod {
		pods, err := k.findPods(r, state.Zone)
		return pods, err
	}
	services, err := k.findServices(r, state.Zone)
	return services, err
}
func serviceFQDN(obj meta.Object, zone string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return dnsutil.Join(obj.GetName(), obj.GetNamespace(), Svc, zone)
}
func podFQDN(p *object.Pod, zone string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if strings.Contains(p.PodIP, ".") {
		name := strings.Replace(p.PodIP, ".", "-", -1)
		return dnsutil.Join(name, p.GetNamespace(), Pod, zone)
	}
	name := strings.Replace(p.PodIP, ":", "-", -1)
	return dnsutil.Join(name, p.GetNamespace(), Pod, zone)
}
func endpointFQDN(ep *object.Endpoints, zone string, endpointNameMode bool) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var names []string
	for _, ss := range ep.Subsets {
		for _, addr := range ss.Addresses {
			names = append(names, dnsutil.Join(endpointHostname(addr, endpointNameMode), serviceFQDN(ep, zone)))
		}
	}
	return names
}
func endpointHostname(addr object.EndpointAddress, endpointNameMode bool) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if addr.Hostname != "" {
		return addr.Hostname
	}
	if endpointNameMode && addr.TargetRefName != "" {
		return addr.TargetRefName
	}
	if strings.Contains(addr.IP, ".") {
		return strings.Replace(addr.IP, ".", "-", -1)
	}
	if strings.Contains(addr.IP, ":") {
		return strings.Replace(addr.IP, ":", "-", -1)
	}
	return ""
}
func (k *Kubernetes) findPods(r recordRequest, zone string) (pods []msg.Service, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if k.podMode == podModeDisabled {
		return nil, errNoItems
	}
	namespace := r.namespace
	podname := r.service
	zonePath := msg.Path(zone, coredns)
	ip := ""
	if podname == "" {
		if k.namespace(namespace) || wildcard(namespace) {
			return nil, nil
		}
		return nil, errNoItems
	}
	if strings.Count(podname, "-") == 3 && !strings.Contains(podname, "--") {
		ip = strings.Replace(podname, "-", ".", -1)
	} else {
		ip = strings.Replace(podname, "-", ":", -1)
	}
	if k.podMode == podModeInsecure {
		if !wildcard(namespace) && !k.namespace(namespace) {
			return nil, errNoItems
		}
		if net.ParseIP(ip) == nil {
			return nil, errNoItems
		}
		return []msg.Service{{Key: strings.Join([]string{zonePath, Pod, namespace, podname}, "/"), Host: ip, TTL: k.ttl}}, err
	}
	err = errNoItems
	if wildcard(podname) && !wildcard(namespace) {
		if k.namespace(namespace) {
			err = nil
		}
	}
	for _, p := range k.APIConn.PodIndex(ip) {
		if wildcard(namespace) && !k.namespaceExposed(p.Namespace) {
			continue
		}
		if p.Deleting {
			continue
		}
		if ip == p.PodIP && match(namespace, p.Namespace) {
			s := msg.Service{Key: strings.Join([]string{zonePath, Pod, namespace, podname}, "/"), Host: ip, TTL: k.ttl}
			pods = append(pods, s)
			err = nil
		}
	}
	return pods, err
}
func (k *Kubernetes) findServices(r recordRequest, zone string) (services []msg.Service, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	zonePath := msg.Path(zone, coredns)
	err = errNoItems
	if wildcard(r.service) && !wildcard(r.namespace) {
		if k.namespace(r.namespace) {
			err = nil
		}
	}
	var (
		endpointsListFunc	func() []*object.Endpoints
		endpointsList		[]*object.Endpoints
		serviceList		[]*object.Service
	)
	if r.service == "" {
		if k.namespace(r.namespace) || wildcard(r.namespace) {
			return nil, nil
		}
		return nil, errNoItems
	}
	if wildcard(r.service) || wildcard(r.namespace) {
		serviceList = k.APIConn.ServiceList()
		endpointsListFunc = func() []*object.Endpoints {
			return k.APIConn.EndpointsList()
		}
	} else {
		idx := object.ServiceKey(r.service, r.namespace)
		serviceList = k.APIConn.SvcIndex(idx)
		endpointsListFunc = func() []*object.Endpoints {
			return k.APIConn.EpIndex(idx)
		}
	}
	for _, svc := range serviceList {
		if !(match(r.namespace, svc.Namespace) && match(r.service, svc.Name)) {
			continue
		}
		if wildcard(r.namespace) && !k.namespaceExposed(svc.Namespace) {
			continue
		}
		if k.opts.ignoreEmptyService && svc.ClusterIP != api.ClusterIPNone {
			podsCount := 0
			for _, ep := range endpointsListFunc() {
				for _, eps := range ep.Subsets {
					podsCount = podsCount + len(eps.Addresses)
				}
			}
			if podsCount == 0 {
				continue
			}
		}
		if svc.ClusterIP == api.ClusterIPNone || r.endpoint != "" {
			if endpointsList == nil {
				endpointsList = endpointsListFunc()
			}
			for _, ep := range endpointsList {
				if ep.Name != svc.Name || ep.Namespace != svc.Namespace {
					continue
				}
				for _, eps := range ep.Subsets {
					for _, addr := range eps.Addresses {
						if r.endpoint != "" {
							if !match(r.endpoint, endpointHostname(addr, k.endpointNameMode)) {
								continue
							}
						}
						for _, p := range eps.Ports {
							if !(match(r.port, p.Name) && match(r.protocol, string(p.Protocol))) {
								continue
							}
							s := msg.Service{Host: addr.IP, Port: int(p.Port), TTL: k.ttl}
							s.Key = strings.Join([]string{zonePath, Svc, svc.Namespace, svc.Name, endpointHostname(addr, k.endpointNameMode)}, "/")
							err = nil
							services = append(services, s)
						}
					}
				}
			}
			continue
		}
		if svc.Type == api.ServiceTypeExternalName {
			s := msg.Service{Key: strings.Join([]string{zonePath, Svc, svc.Namespace, svc.Name}, "/"), Host: svc.ExternalName, TTL: k.ttl}
			if t, _ := s.HostType(); t == dns.TypeCNAME {
				s.Key = strings.Join([]string{zonePath, Svc, svc.Namespace, svc.Name}, "/")
				services = append(services, s)
				err = nil
			}
			continue
		}
		for _, p := range svc.Ports {
			if !(match(r.port, p.Name) && match(r.protocol, string(p.Protocol))) {
				continue
			}
			err = nil
			s := msg.Service{Host: svc.ClusterIP, Port: int(p.Port), TTL: k.ttl}
			s.Key = strings.Join([]string{zonePath, Svc, svc.Namespace, svc.Name}, "/")
			services = append(services, s)
		}
	}
	return services, err
}
func match(a, b string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if wildcard(a) {
		return true
	}
	if wildcard(b) {
		return true
	}
	return strings.EqualFold(a, b)
}
func wildcard(s string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return s == "*" || s == "any"
}

const coredns = "c"
