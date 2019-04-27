package kubernetes

import (
	"testing"
	"github.com/mholt/caddy"
)

func TestKubernetesParseTransfer(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tests := []struct {
		input		string
		expected	string
		shouldErr	bool
	}{{`kubernetes cluster.local {
			transfer to 1.2.3.4
		}`, "1.2.3.4:53", false}, {`kubernetes cluster.local {
			transfer to 1.2.3.4:53
		}`, "1.2.3.4:53", false}, {`kubernetes cluster.local {
			transfer to *
		}`, "*", false}, {`kubernetes cluster.local {
			transfer
		}`, "", true}}
	for i, tc := range tests {
		c := caddy.NewTestController("dns", tc.input)
		k, err := kubernetesParse(c)
		if err != nil && !tc.shouldErr {
			t.Fatalf("Test %d: Expected no error, got %q", i, err)
		}
		if err == nil && tc.shouldErr {
			t.Fatalf("Test %d: Expected error, got none", i)
		}
		if err != nil && tc.shouldErr {
			continue
		}
		if k.TransferTo[0] != tc.expected {
			t.Errorf("Test %d: Expected Transfer To to be %s, got %s", i, tc.expected, k.TransferTo[0])
		}
	}
}
