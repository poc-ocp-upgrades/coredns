package etcd

import (
	"net"
	"strconv"
	"time"
	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/plugin/pkg/dnsutil"
	"github.com/coredns/coredns/plugin/proxy"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func (e *Etcd) UpdateStubZones() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	go func() {
		for {
			e.updateStubZones()
			time.Sleep(15 * time.Second)
		}
	}()
}
func (e *Etcd) updateStubZones() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	zone := e.Zones[0]
	fakeState := request.Request{W: nil, Req: new(dns.Msg)}
	fakeState.Req.SetQuestion(stubDomain+"."+zone, dns.TypeA)
	services, err := e.Records(fakeState, false)
	if err != nil {
		return
	}
	stubmap := make(map[string]proxy.Proxy)
	nameservers := map[string][]string{}
Services:
	for _, serv := range services {
		if serv.Port == 0 {
			serv.Port = 53
		}
		ip := net.ParseIP(serv.Host)
		if ip == nil {
			log.Warningf("Non IP address stub nameserver: %s", serv.Host)
			continue
		}
		domain := msg.Domain(serv.Key)
		labels := dns.SplitDomainName(domain)
		for _, z := range e.Zones {
			domain = dnsutil.Join(labels[1 : len(labels)-dns.CountLabel(z)-2]...)
			if domain == z {
				log.Warningf("Skipping nameserver for domain we are authoritative for: %s", domain)
				continue Services
			}
		}
		nameservers[domain] = append(nameservers[domain], net.JoinHostPort(serv.Host, strconv.Itoa(serv.Port)))
	}
	for domain, nss := range nameservers {
		stubmap[domain] = proxy.NewLookup(nss)
	}
	if len(stubmap) > 0 {
		e.Stubmap = &stubmap
	}
	return
}
