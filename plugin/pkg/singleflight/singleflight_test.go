package singleflight

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var g Group
	v, err := g.Do(1, func() (interface{}, error) {
		return "bar", nil
	})
	if got, want := fmt.Sprintf("%v (%T)", v, v), "bar (string)"; got != want {
		t.Errorf("Do = %v; want %v", got, want)
	}
	if err != nil {
		t.Errorf("Do error = %v", err)
	}
}
func TestDoErr(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var g Group
	someErr := errors.New("Some error")
	v, err := g.Do(1, func() (interface{}, error) {
		return nil, someErr
	})
	if err != someErr {
		t.Errorf("Do error = %v; want someErr", err)
	}
	if v != nil {
		t.Errorf("Unexpected non-nil value %#v", v)
	}
}
func TestDoDupSuppress(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var g Group
	c := make(chan string)
	var calls int32
	fn := func() (interface{}, error) {
		atomic.AddInt32(&calls, 1)
		return <-c, nil
	}
	const n = 10
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			v, err := g.Do(1, fn)
			if err != nil {
				t.Errorf("Do error: %v", err)
			}
			if v.(string) != "bar" {
				t.Errorf("Got %q; want %q", v, "bar")
			}
			wg.Done()
		}()
	}
	time.Sleep(100 * time.Millisecond)
	c <- "bar"
	wg.Wait()
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Errorf("Number of calls = %d; want 1", got)
	}
}
