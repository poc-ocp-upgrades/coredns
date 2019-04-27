package erratic

import (
	"sync/atomic"
)

func (e *Erratic) Health() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	q := atomic.LoadUint64(&e.q)
	if e.drop > 0 && q%e.drop == 0 {
		return false
	}
	return true
}
