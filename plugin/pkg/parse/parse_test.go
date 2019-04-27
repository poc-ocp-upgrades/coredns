package parse

import (
	"testing"
	"github.com/mholt/caddy"
)

func TestTransfer(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	tests := []struct {
		inputFileRules	string
		shouldErr	bool
		secondary	bool
		expectedTo	[]string
		expectedFrom	[]string
	}{{`to 127.0.0.1`, false, false, []string{"127.0.0.1:53"}, []string{}}, {`to 127.0.0.1 127.0.0.2`, false, false, []string{"127.0.0.1:53", "127.0.0.2:53"}, []string{}}, {`from 127.0.0.1`, false, true, []string{}, []string{"127.0.0.1:53"}}, {`from 127.0.0.1 127.0.0.2`, false, true, []string{}, []string{"127.0.0.1:53", "127.0.0.2:53"}}, {`to 127.0.0.1 127.0.0.2
			from 127.0.0.1 127.0.0.2`, false, true, []string{"127.0.0.1:53", "127.0.0.2:53"}, []string{"127.0.0.1:53", "127.0.0.2:53"}}, {`from 127.0.0.1`, true, false, []string{}, []string{}}, {`from !@#$%^&*()`, true, true, []string{}, []string{}}, {`from`, true, false, []string{}, []string{}}, {`from *`, true, true, []string{}, []string{}}}
	for i, test := range tests {
		c := caddy.NewTestController("dns", test.inputFileRules)
		tos, froms, err := Transfer(c, test.secondary)
		if err == nil && test.shouldErr {
			t.Fatalf("Test %d expected errors, but got no error %+v %+v", i, err, test)
		} else if err != nil && !test.shouldErr {
			t.Fatalf("Test %d expected no errors, but got '%v'", i, err)
		}
		if test.expectedTo != nil {
			for j, got := range tos {
				if got != test.expectedTo[j] {
					t.Fatalf("Test %d expected %v, got %v", i, test.expectedTo[j], got)
				}
			}
		}
		if test.expectedFrom != nil {
			for j, got := range froms {
				if got != test.expectedFrom[j] {
					t.Fatalf("Test %d expected %v, got %v", i, test.expectedFrom[j], got)
				}
			}
		}
	}
}
