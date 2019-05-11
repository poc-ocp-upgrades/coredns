package kubernetes

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"github.com/coredns/coredns/plugin/kubernetes/object"
	dnswatch "github.com/coredns/coredns/plugin/pkg/watch"
	api "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

const (
	podIPIndex				= "PodIP"
	svcNameNamespaceIndex	= "NameNamespace"
	svcIPIndex				= "ServiceIP"
	epNameNamespaceIndex	= "EndpointNameNamespace"
	epIPIndex				= "EndpointsIP"
)

type dnsController interface {
	ServiceList() []*object.Service
	EndpointsList() []*object.Endpoints
	SvcIndex(string) []*object.Service
	SvcIndexReverse(string) []*object.Service
	PodIndex(string) []*object.Pod
	EpIndex(string) []*object.Endpoints
	EpIndexReverse(string) []*object.Endpoints
	GetNodeByName(string) (*api.Node, error)
	GetNamespaceByName(string) (*api.Namespace, error)
	Run()
	HasSynced() bool
	Stop() error
	Modified() int64
	SetWatchChan(dnswatch.Chan)
	Watch(string) error
	StopWatching(string)
}
type dnsControl struct {
	modified			int64
	client				kubernetes.Interface
	selector			labels.Selector
	svcController		cache.Controller
	podController		cache.Controller
	epController		cache.Controller
	nsController		cache.Controller
	svcLister			cache.Indexer
	podLister			cache.Indexer
	epLister			cache.Indexer
	nsLister			cache.Store
	stopLock			sync.Mutex
	shutdown			bool
	stopCh				chan struct{}
	watchChan			dnswatch.Chan
	watched				map[string]struct{}
	zones				[]string
	endpointNameMode	bool
}
type dnsControlOpts struct {
	initPodCache		bool
	initEndpointsCache	bool
	resyncPeriod		time.Duration
	ignoreEmptyService	bool
	labelSelector		*meta.LabelSelector
	selector			labels.Selector
	zones				[]string
	endpointNameMode	bool
}

func newdnsController(kubeClient kubernetes.Interface, opts dnsControlOpts) *dnsControl {
	_logClusterCodePath()
	defer _logClusterCodePath()
	dns := dnsControl{client: kubeClient, selector: opts.selector, stopCh: make(chan struct{}), watched: make(map[string]struct{}), zones: opts.zones, endpointNameMode: opts.endpointNameMode}
	dns.svcLister, dns.svcController = object.NewIndexerInformer(&cache.ListWatch{ListFunc: serviceListFunc(dns.client, api.NamespaceAll, dns.selector), WatchFunc: serviceWatchFunc(dns.client, api.NamespaceAll, dns.selector)}, &api.Service{}, opts.resyncPeriod, cache.ResourceEventHandlerFuncs{AddFunc: dns.Add, UpdateFunc: dns.Update, DeleteFunc: dns.Delete}, cache.Indexers{svcNameNamespaceIndex: svcNameNamespaceIndexFunc, svcIPIndex: svcIPIndexFunc}, object.ToService)
	if opts.initPodCache {
		dns.podLister, dns.podController = object.NewIndexerInformer(&cache.ListWatch{ListFunc: podListFunc(dns.client, api.NamespaceAll, dns.selector), WatchFunc: podWatchFunc(dns.client, api.NamespaceAll, dns.selector)}, &api.Pod{}, opts.resyncPeriod, cache.ResourceEventHandlerFuncs{AddFunc: dns.Add, UpdateFunc: dns.Update, DeleteFunc: dns.Delete}, cache.Indexers{podIPIndex: podIPIndexFunc}, object.ToPod)
	}
	if opts.initEndpointsCache {
		dns.epLister, dns.epController = object.NewIndexerInformer(&cache.ListWatch{ListFunc: endpointsListFunc(dns.client, api.NamespaceAll, dns.selector), WatchFunc: endpointsWatchFunc(dns.client, api.NamespaceAll, dns.selector)}, &api.Endpoints{}, opts.resyncPeriod, cache.ResourceEventHandlerFuncs{AddFunc: dns.Add, UpdateFunc: dns.Update, DeleteFunc: dns.Delete}, cache.Indexers{epNameNamespaceIndex: epNameNamespaceIndexFunc, epIPIndex: epIPIndexFunc}, object.ToEndpoints)
	}
	dns.nsLister, dns.nsController = cache.NewInformer(&cache.ListWatch{ListFunc: namespaceListFunc(dns.client, dns.selector), WatchFunc: namespaceWatchFunc(dns.client, dns.selector)}, &api.Namespace{}, opts.resyncPeriod, cache.ResourceEventHandlerFuncs{})
	return &dns
}
func podIPIndexFunc(obj interface{}) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	p, ok := obj.(*object.Pod)
	if !ok {
		return nil, errObj
	}
	return []string{p.PodIP}, nil
}
func svcIPIndexFunc(obj interface{}) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	svc, ok := obj.(*object.Service)
	if !ok {
		return nil, errObj
	}
	if len(svc.ExternalIPs) == 0 {
		return []string{svc.ClusterIP}, nil
	}
	return append([]string{svc.ClusterIP}, svc.ExternalIPs...), nil
}
func svcNameNamespaceIndexFunc(obj interface{}) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s, ok := obj.(*object.Service)
	if !ok {
		return nil, errObj
	}
	return []string{s.Index}, nil
}
func epNameNamespaceIndexFunc(obj interface{}) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s, ok := obj.(*object.Endpoints)
	if !ok {
		return nil, errObj
	}
	return []string{s.Index}, nil
}
func epIPIndexFunc(obj interface{}) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ep, ok := obj.(*object.Endpoints)
	if !ok {
		return nil, errObj
	}
	return ep.IndexIP, nil
}
func serviceListFunc(c kubernetes.Interface, ns string, s labels.Selector) func(meta.ListOptions) (runtime.Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(opts meta.ListOptions) (runtime.Object, error) {
		if s != nil {
			opts.LabelSelector = s.String()
		}
		listV1, err := c.CoreV1().Services(ns).List(opts)
		return listV1, err
	}
}
func podListFunc(c kubernetes.Interface, ns string, s labels.Selector) func(meta.ListOptions) (runtime.Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(opts meta.ListOptions) (runtime.Object, error) {
		if s != nil {
			opts.LabelSelector = s.String()
		}
		listV1, err := c.CoreV1().Pods(ns).List(opts)
		return listV1, err
	}
}
func serviceWatchFunc(c kubernetes.Interface, ns string, s labels.Selector) func(options meta.ListOptions) (watch.Interface, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(options meta.ListOptions) (watch.Interface, error) {
		if s != nil {
			options.LabelSelector = s.String()
		}
		w, err := c.CoreV1().Services(ns).Watch(options)
		return w, err
	}
}
func podWatchFunc(c kubernetes.Interface, ns string, s labels.Selector) func(options meta.ListOptions) (watch.Interface, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(options meta.ListOptions) (watch.Interface, error) {
		if s != nil {
			options.LabelSelector = s.String()
		}
		w, err := c.CoreV1().Pods(ns).Watch(options)
		return w, err
	}
}
func endpointsListFunc(c kubernetes.Interface, ns string, s labels.Selector) func(meta.ListOptions) (runtime.Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(opts meta.ListOptions) (runtime.Object, error) {
		if s != nil {
			opts.LabelSelector = s.String()
		}
		listV1, err := c.CoreV1().Endpoints(ns).List(opts)
		return listV1, err
	}
}
func endpointsWatchFunc(c kubernetes.Interface, ns string, s labels.Selector) func(options meta.ListOptions) (watch.Interface, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(options meta.ListOptions) (watch.Interface, error) {
		if s != nil {
			options.LabelSelector = s.String()
		}
		w, err := c.CoreV1().Endpoints(ns).Watch(options)
		return w, err
	}
}
func namespaceListFunc(c kubernetes.Interface, s labels.Selector) func(meta.ListOptions) (runtime.Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(opts meta.ListOptions) (runtime.Object, error) {
		if s != nil {
			opts.LabelSelector = s.String()
		}
		listV1, err := c.CoreV1().Namespaces().List(opts)
		return listV1, err
	}
}
func namespaceWatchFunc(c kubernetes.Interface, s labels.Selector) func(options meta.ListOptions) (watch.Interface, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(options meta.ListOptions) (watch.Interface, error) {
		if s != nil {
			options.LabelSelector = s.String()
		}
		w, err := c.CoreV1().Namespaces().Watch(options)
		return w, err
	}
}
func (dns *dnsControl) SetWatchChan(c dnswatch.Chan) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	dns.watchChan = c
}
func (dns *dnsControl) StopWatching(qname string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	delete(dns.watched, qname)
}
func (dns *dnsControl) Watch(qname string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if dns.watchChan == nil {
		return fmt.Errorf("cannot start watch because the channel has not been set")
	}
	dns.watched[qname] = struct{}{}
	return nil
}
func (dns *dnsControl) Stop() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	dns.stopLock.Lock()
	defer dns.stopLock.Unlock()
	if !dns.shutdown {
		close(dns.stopCh)
		dns.shutdown = true
		return nil
	}
	return fmt.Errorf("shutdown already in progress")
}
func (dns *dnsControl) Run() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	go dns.svcController.Run(dns.stopCh)
	if dns.epController != nil {
		go dns.epController.Run(dns.stopCh)
	}
	if dns.podController != nil {
		go dns.podController.Run(dns.stopCh)
	}
	go dns.nsController.Run(dns.stopCh)
	<-dns.stopCh
}
func (dns *dnsControl) HasSynced() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	a := dns.svcController.HasSynced()
	b := true
	if dns.epController != nil {
		b = dns.epController.HasSynced()
	}
	c := true
	if dns.podController != nil {
		c = dns.podController.HasSynced()
	}
	d := dns.nsController.HasSynced()
	return a && b && c && d
}
func (dns *dnsControl) ServiceList() (svcs []*object.Service) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	os := dns.svcLister.List()
	for _, o := range os {
		s, ok := o.(*object.Service)
		if !ok {
			continue
		}
		svcs = append(svcs, s)
	}
	return svcs
}
func (dns *dnsControl) EndpointsList() (eps []*object.Endpoints) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	os := dns.epLister.List()
	for _, o := range os {
		ep, ok := o.(*object.Endpoints)
		if !ok {
			continue
		}
		eps = append(eps, ep)
	}
	return eps
}
func (dns *dnsControl) PodIndex(ip string) (pods []*object.Pod) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	os, err := dns.podLister.ByIndex(podIPIndex, ip)
	if err != nil {
		return nil
	}
	for _, o := range os {
		p, ok := o.(*object.Pod)
		if !ok {
			continue
		}
		pods = append(pods, p)
	}
	return pods
}
func (dns *dnsControl) SvcIndex(idx string) (svcs []*object.Service) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	os, err := dns.svcLister.ByIndex(svcNameNamespaceIndex, idx)
	if err != nil {
		return nil
	}
	for _, o := range os {
		s, ok := o.(*object.Service)
		if !ok {
			continue
		}
		svcs = append(svcs, s)
	}
	return svcs
}
func (dns *dnsControl) SvcIndexReverse(ip string) (svcs []*object.Service) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	os, err := dns.svcLister.ByIndex(svcIPIndex, ip)
	if err != nil {
		return nil
	}
	for _, o := range os {
		s, ok := o.(*object.Service)
		if !ok {
			continue
		}
		svcs = append(svcs, s)
	}
	return svcs
}
func (dns *dnsControl) EpIndex(idx string) (ep []*object.Endpoints) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	os, err := dns.epLister.ByIndex(epNameNamespaceIndex, idx)
	if err != nil {
		return nil
	}
	for _, o := range os {
		e, ok := o.(*object.Endpoints)
		if !ok {
			continue
		}
		ep = append(ep, e)
	}
	return ep
}
func (dns *dnsControl) EpIndexReverse(ip string) (ep []*object.Endpoints) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	os, err := dns.epLister.ByIndex(epIPIndex, ip)
	if err != nil {
		return nil
	}
	for _, o := range os {
		e, ok := o.(*object.Endpoints)
		if !ok {
			continue
		}
		ep = append(ep, e)
	}
	return ep
}
func (dns *dnsControl) GetNodeByName(name string) (*api.Node, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	v1node, err := dns.client.CoreV1().Nodes().Get(name, meta.GetOptions{})
	return v1node, err
}
func (dns *dnsControl) GetNamespaceByName(name string) (*api.Namespace, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	os := dns.nsLister.List()
	for _, o := range os {
		ns, ok := o.(*api.Namespace)
		if !ok {
			continue
		}
		if name == ns.ObjectMeta.Name {
			return ns, nil
		}
	}
	return nil, fmt.Errorf("namespace not found")
}
func (dns *dnsControl) Modified() int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	unix := atomic.LoadInt64(&dns.modified)
	return unix
}
func (dns *dnsControl) updateModifed() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	unix := time.Now().Unix()
	atomic.StoreInt64(&dns.modified, unix)
}

var errObj = errors.New("obj was not of the correct type")
