package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var workableServer *httptest.Server

func TestMain(m *testing.M) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	workableServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	r := m.Run()
	workableServer.Close()
	os.Exit(r)
}

type customPolicy struct{}

func (r *customPolicy) Select(pool HostPool) *UpstreamHost {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return pool[0]
}
func testPool() HostPool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pool := []*UpstreamHost{{Name: workableServer.URL}, {Name: "http://shouldnot.resolve:85"}, {Name: "http://C"}}
	return HostPool(pool)
}
func TestRegisterPolicy(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	name := "custom"
	customPolicy := &customPolicy{}
	RegisterPolicy(name, func() Policy {
		return customPolicy
	})
	if _, ok := SupportedPolicies[name]; !ok {
		t.Error("Expected supportedPolicies to have a custom policy.")
	}
}
func TestHealthCheck(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	u := &HealthCheck{Hosts: testPool(), Path: "/", FailTimeout: 10 * time.Second, MaxFails: 1}
	for i, h := range u.Hosts {
		u.Hosts[i].CheckURL = u.normalizeCheckURL(h.Name)
	}
	u.healthCheck()
	time.Sleep(time.Duration(1 * time.Second))
	if u.Hosts[0].Down() {
		t.Error("Expected first host in testpool to not fail healthcheck.")
	}
	if !u.Hosts[1].Down() {
		t.Error("Expected second host in testpool to fail healthcheck.")
	}
}
func TestHealthCheckDisabled(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	u := &HealthCheck{Hosts: testPool(), FailTimeout: 10 * time.Second, MaxFails: 1}
	for i, h := range u.Hosts {
		u.Hosts[i].CheckURL = u.normalizeCheckURL(h.Name)
	}
	u.healthCheck()
	time.Sleep(time.Duration(1 * time.Second))
	for i, h := range u.Hosts {
		if h.Down() {
			t.Errorf("Expected host %d in testpool to not be down with healthchecks disabled.", i+1)
		}
	}
}
func TestRoundRobinPolicy(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pool := testPool()
	rrPolicy := &RoundRobin{}
	h := rrPolicy.Select(pool)
	if h != pool[1] {
		t.Error("Expected first round robin host to be second host in the pool.")
	}
	h = rrPolicy.Select(pool)
	if h != pool[2] {
		t.Error("Expected second round robin host to be third host in the pool.")
	}
	h = rrPolicy.Select(pool)
	if h != pool[0] {
		t.Error("Expected third round robin host to be first host in the pool.")
	}
}
func TestLeastConnPolicy(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pool := testPool()
	lcPolicy := &LeastConn{}
	pool[0].Conns = 10
	pool[1].Conns = 10
	h := lcPolicy.Select(pool)
	if h != pool[2] {
		t.Error("Expected least connection host to be third host.")
	}
	pool[2].Conns = 100
	h = lcPolicy.Select(pool)
	if h != pool[0] && h != pool[1] {
		t.Error("Expected least connection host to be first or second host.")
	}
}
func TestCustomPolicy(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pool := testPool()
	customPolicy := &customPolicy{}
	h := customPolicy.Select(pool)
	if h != pool[0] {
		t.Error("Expected custom policy host to be the first host.")
	}
}
func TestFirstPolicy(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pool := testPool()
	rrPolicy := &First{}
	h := rrPolicy.Select(pool)
	if h != pool[0] {
		t.Error("Expected always first to be first host in the pool.")
	}
	h = rrPolicy.Select(pool)
	if h != pool[0] {
		t.Error("Expected always first to be first host in the pool, even in second call")
	}
	pool[0].Fails = 1
	h = rrPolicy.Select(pool)
	if h != pool[1] {
		t.Error("Expected first to be he second in pool if the first one is down.")
	}
}
