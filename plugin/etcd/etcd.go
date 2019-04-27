package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/plugin/pkg/fall"
	"github.com/coredns/coredns/plugin/proxy"
	"github.com/coredns/coredns/request"
	"github.com/coredns/coredns/plugin/pkg/upstream"
	etcdcv3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/miekg/dns"
)

const (
	priority	= 10
	ttl		= 300
	etcdTimeout	= 5 * time.Second
)

var errKeyNotFound = errors.New("Key not found")

type Etcd struct {
	Next		plugin.Handler
	Fall		fall.F
	Zones		[]string
	PathPrefix	string
	Upstream	upstream.Upstream
	Client		*etcdcv3.Client
	Ctx		context.Context
	Stubmap		*map[string]proxy.Proxy
	endpoints	[]string
}

func (e *Etcd) Services(state request.Request, exact bool, opt plugin.Options) (services []msg.Service, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	services, err = e.Records(state, exact)
	if err != nil {
		return
	}
	services = msg.Group(services)
	return
}
func (e *Etcd) Reverse(state request.Request, exact bool, opt plugin.Options) (services []msg.Service, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return e.Services(state, exact, opt)
}
func (e *Etcd) Lookup(state request.Request, name string, typ uint16) (*dns.Msg, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return e.Upstream.Lookup(state, name, typ)
}
func (e *Etcd) IsNameError(err error) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return err == errKeyNotFound
}
func (e *Etcd) Records(state request.Request, exact bool) ([]msg.Service, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	name := state.Name()
	path, star := msg.PathWithWildcard(name, e.PathPrefix)
	r, err := e.get(path, !exact)
	if err != nil {
		return nil, err
	}
	segments := strings.Split(msg.Path(name, e.PathPrefix), "/")
	return e.loopNodes(r.Kvs, segments, star)
}
func (e *Etcd) get(path string, recursive bool) (*etcdcv3.GetResponse, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ctx, cancel := context.WithTimeout(e.Ctx, etcdTimeout)
	defer cancel()
	if recursive == true {
		if !strings.HasSuffix(path, "/") {
			path = path + "/"
		}
		r, err := e.Client.Get(ctx, path, etcdcv3.WithPrefix())
		if err != nil {
			return nil, err
		}
		if r.Count == 0 {
			path = strings.TrimSuffix(path, "/")
			r, err = e.Client.Get(ctx, path)
			if err != nil {
				return nil, err
			}
			if r.Count == 0 {
				return nil, errKeyNotFound
			}
		}
		return r, nil
	}
	r, err := e.Client.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if r.Count == 0 {
		return nil, errKeyNotFound
	}
	return r, nil
}
func (e *Etcd) loopNodes(kv []*mvccpb.KeyValue, nameParts []string, star bool) (sx []msg.Service, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	bx := make(map[msg.Service]struct{})
Nodes:
	for _, n := range kv {
		if star {
			s := string(n.Key)
			keyParts := strings.Split(s, "/")
			for i, n := range nameParts {
				if i > len(keyParts)-1 {
					continue Nodes
				}
				if n == "*" || n == "any" {
					continue
				}
				if keyParts[i] != n {
					continue Nodes
				}
			}
		}
		serv := new(msg.Service)
		if err := json.Unmarshal(n.Value, serv); err != nil {
			return nil, fmt.Errorf("%s: %s", n.Key, err.Error())
		}
		b := msg.Service{Host: serv.Host, Port: serv.Port, Priority: serv.Priority, Weight: serv.Weight, Text: serv.Text, Key: string(n.Key)}
		if _, ok := bx[b]; ok {
			continue
		}
		bx[b] = struct{}{}
		serv.Key = string(n.Key)
		serv.TTL = e.TTL(n, serv)
		if serv.Priority == 0 {
			serv.Priority = priority
		}
		sx = append(sx, *serv)
	}
	return sx, nil
}
func (e *Etcd) TTL(kv *mvccpb.KeyValue, serv *msg.Service) uint32 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	etcdTTL := uint32(kv.Lease)
	if etcdTTL == 0 && serv.TTL == 0 {
		return ttl
	}
	if etcdTTL == 0 {
		return serv.TTL
	}
	if serv.TTL == 0 {
		return etcdTTL
	}
	if etcdTTL < serv.TTL {
		return etcdTTL
	}
	return serv.TTL
}
