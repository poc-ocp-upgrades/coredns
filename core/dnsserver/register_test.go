package dnsserver

import (
	"testing"
)

func TestHandler(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	tp := testPlugin{}
	c := testConfig("dns", tp)
	if _, err := NewServer("127.0.0.1:53", []*Config{c}); err != nil {
		t.Errorf("Expected no error for NewServer, got %s", err)
	}
	if h := c.Handler("testplugin"); h != tp {
		t.Errorf("Expected testPlugin from Handler, got %T", h)
	}
	if h := c.Handler("nothing"); h != nil {
		t.Errorf("Expected nil from Handler, got %T", h)
	}
}
func TestHandlers(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	tp := testPlugin{}
	c := testConfig("dns", tp)
	if _, err := NewServer("127.0.0.1:53", []*Config{c}); err != nil {
		t.Errorf("Expected no error for NewServer, got %s", err)
	}
	hs := c.Handlers()
	if len(hs) != 1 || hs[0] != tp {
		t.Errorf("Expected [testPlugin] from Handlers, got %v", hs)
	}
}
func TestGroupingServers(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for i, test := range []struct {
		configs		[]*Config
		expectedGroups	[]string
		failing		bool
	}{{configs: []*Config{{Transport: "dns", Zone: ".", Port: "53", ListenHosts: []string{""}}}, expectedGroups: []string{"dns://:53"}, failing: false}, {configs: []*Config{{Transport: "dns", Zone: ".", Port: "53", ListenHosts: []string{""}}, {Transport: "dns", Zone: ".", Port: "54", ListenHosts: []string{""}}}, expectedGroups: []string{"dns://:53", "dns://:54"}, failing: false}, {configs: []*Config{{Transport: "dns", Zone: ".", Port: "53", ListenHosts: []string{""}}, {Transport: "dns", Zone: "com.", Port: "53", ListenHosts: []string{""}}}, expectedGroups: []string{"dns://:53"}, failing: false}, {configs: []*Config{{Transport: "dns", Zone: ".", Port: "53", ListenHosts: []string{"127.0.0.1"}}, {Transport: "dns", Zone: ".", Port: "54", ListenHosts: []string{""}}}, expectedGroups: []string{"dns://127.0.0.1:53", "dns://:54"}, failing: false}, {configs: []*Config{{Transport: "dns", Zone: ".", Port: "53", ListenHosts: []string{"127.0.0.1", "::1"}}, {Transport: "dns", Zone: ".", Port: "54", ListenHosts: []string{""}}}, expectedGroups: []string{"dns://127.0.0.1:53", "dns://[::1]:53", "dns://:54"}, failing: false}, {configs: []*Config{{Transport: "dns", Zone: ".", Port: "53", ListenHosts: []string{"127.0.0.1", "::1"}}, {Transport: "dns", Zone: "com.", Port: "53", ListenHosts: []string{"127.0.0.1", "::1"}}}, expectedGroups: []string{"dns://127.0.0.1:53", "dns://[::1]:53"}, failing: false}, {configs: []*Config{{Transport: "dns", Zone: ".", Port: "53", ListenHosts: []string{"127.0.0.1"}}, {Transport: "dns", Zone: "com.", Port: "53", ListenHosts: []string{"::1"}}}, expectedGroups: []string{"dns://127.0.0.1:53", "dns://[::1]:53"}, failing: false}, {configs: []*Config{{Transport: "dns", Zone: ".", Port: "53", ListenHosts: []string{"127.0.0.1", "::1"}}, {Transport: "dns", Zone: "com.", Port: "53", ListenHosts: []string{""}}}, expectedGroups: []string{"dns://127.0.0.1:53", "dns://[::1]:53", "dns://:53"}, failing: false}} {
		groups, err := groupConfigsByListenAddr(test.configs)
		if err != nil {
			if !test.failing {
				t.Fatalf("Test %d, expected no errors, but got: %v", i, err)
			}
			continue
		}
		if test.failing {
			t.Fatalf("Test %d, expected to failed but did not, returned values", i)
		}
		if len(groups) != len(test.expectedGroups) {
			t.Errorf("Test %d : expected the group's size to be %d, was %d", i, len(test.expectedGroups), len(groups))
			continue
		}
		for _, v := range test.expectedGroups {
			if _, ok := groups[v]; !ok {
				t.Errorf("Test %d : expected value %v to be in the group, was not", i, v)
			}
		}
	}
}
