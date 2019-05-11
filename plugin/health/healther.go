package health

type Healther interface{ Health() bool }

func (h *health) Ok() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	h.RLock()
	defer h.RUnlock()
	return h.ok
}
func (h *health) SetOk(ok bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	h.Lock()
	defer h.Unlock()
	h.ok = ok
}
func (h *health) poll() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, m := range h.h {
		if !m.Health() {
			h.SetOk(false)
			return
		}
	}
	h.SetOk(true)
}
