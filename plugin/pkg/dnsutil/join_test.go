package dnsutil

import "testing"

func TestJoin(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	tests := []struct {
		in	[]string
		out	string
	}{{[]string{"bla", "bliep", "example", "org"}, "bla.bliep.example.org."}, {[]string{"example", "."}, "example."}, {[]string{"example", "org."}, "example.org."}, {[]string{"."}, "."}}
	for i, tc := range tests {
		if x := Join(tc.in...); x != tc.out {
			t.Errorf("Test %d, expected %s, got %s", i, tc.out, x)
		}
	}
}
