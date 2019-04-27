package freq

import (
	"testing"
	"time"
)

func TestFreqUpdate(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	now := time.Now().UTC()
	f := New(now)
	window := 1 * time.Minute
	f.Update(window, time.Now().UTC())
	f.Update(window, time.Now().UTC())
	f.Update(window, time.Now().UTC())
	hitsCheck(t, f, 3)
	f.Reset(now, 0)
	history := time.Now().UTC().Add(-3 * time.Minute)
	f.Update(window, history)
	hitsCheck(t, f, 1)
}
func TestReset(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	f := New(time.Now().UTC())
	f.Update(1*time.Minute, time.Now().UTC())
	hitsCheck(t, f, 1)
	f.Reset(time.Now().UTC(), 0)
	hitsCheck(t, f, 0)
}
func hitsCheck(t *testing.T, f *Freq, expected int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if x := f.Hits(); x != expected {
		t.Fatalf("Expected hits to be %d, got %d", expected, x)
	}
}
