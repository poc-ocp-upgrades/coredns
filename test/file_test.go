package test

import "testing"

func TestTempFile(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.Parallel()
	_, f, e := TempFile(".", "test")
	if e != nil {
		t.Fatalf("Failed to create temp file: %s", e)
	}
	defer f()
}
