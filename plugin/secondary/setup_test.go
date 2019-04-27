package secondary

import (
	"testing"
	"github.com/mholt/caddy"
)

func TestSecondaryParse(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	tests := []struct {
		inputFileRules	string
		shouldErr	bool
		transferFrom	string
		zones		[]string
	}{{`secondary`, false, "", nil}, {`secondary {
				transfer from 127.0.0.1
				transfer to 127.0.0.1
			}`, false, "127.0.0.1:53", nil}, {`secondary example.org {
				transfer from 127.0.0.1
				transfer to 127.0.0.1
			}`, false, "127.0.0.1:53", []string{"example.org."}}}
	for i, test := range tests {
		c := caddy.NewTestController("dns", test.inputFileRules)
		s, err := secondaryParse(c)
		if err == nil && test.shouldErr {
			t.Fatalf("Test %d expected errors, but got no error", i)
		} else if err != nil && !test.shouldErr {
			t.Fatalf("Test %d expected no errors, but got '%v'", i, err)
		}
		for i, name := range test.zones {
			if x := s.Names[i]; x != name {
				t.Fatalf("Test %d zone names don't match expected %q, but got %q", i, name, x)
			}
		}
		for _, v := range s.Z {
			if x := v.TransferFrom[0]; x != test.transferFrom {
				t.Fatalf("Test %d transform from names don't match expected %q, but got %q", i, test.transferFrom, x)
			}
		}
	}
}
