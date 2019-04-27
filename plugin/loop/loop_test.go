package loop

import "testing"

func TestLoop(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := New(".")
	l.inc()
	if l.seen() != 1 {
		t.Errorf("Failed to inc loop, expected %d, got %d", 1, l.seen())
	}
}
