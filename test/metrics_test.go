package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
	"github.com/coredns/coredns/plugin/cache"
	"github.com/coredns/coredns/plugin/metrics"
	mtest "github.com/coredns/coredns/plugin/metrics/test"
	"github.com/coredns/coredns/plugin/metrics/vars"
	"github.com/miekg/dns"
)

func TestMetricsServer(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	corefile := `example.org:0 {
	chaos CoreDNS-001 miek@miek.nl
	prometheus localhost:0
}

example.com:0 {
	proxy . 8.8.4.4:53
	prometheus localhost:0
}
`
	srv, err := CoreDNSServer(corefile)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}
	defer srv.Stop()
}
func TestMetricsRefused(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	metricName := "coredns_dns_response_rcode_count_total"
	corefile := `example.org:0 {
	proxy . 8.8.8.8:53
	prometheus localhost:0
}
`
	srv, udp, _, err := CoreDNSServerAndPorts(corefile)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}
	defer srv.Stop()
	m := new(dns.Msg)
	m.SetQuestion("google.com.", dns.TypeA)
	if _, err = dns.Exchange(m, udp); err != nil {
		t.Fatalf("Could not send message: %s", err)
	}
	data := mtest.Scrape(t, "http://"+metrics.ListenAddr+"/metrics")
	got, labels := mtest.MetricValue(metricName, data)
	if got != "1" {
		t.Errorf("Expected value %s for refused, but got %s", "1", got)
	}
	if labels["zone"] != vars.Dropped {
		t.Errorf("Expected zone value %s for refused, but got %s", vars.Dropped, labels["zone"])
	}
	if labels["rcode"] != "REFUSED" {
		t.Errorf("Expected zone value %s for refused, but got %s", "REFUSED", labels["rcode"])
	}
}
func TestMetricsCache(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cacheSizeMetricName := "coredns_cache_size"
	cacheHitMetricName := "coredns_cache_hits_total"
	corefile := `example.net:0 {
	proxy . 8.8.8.8:53
	prometheus localhost:0
	cache
}
`
	srv, err := CoreDNSServer(corefile)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}
	defer srv.Stop()
	udp, _ := CoreDNSServerPorts(srv, 0)
	m := new(dns.Msg)
	m.SetQuestion("example.net.", dns.TypeA)
	if _, err = dns.Exchange(m, udp); err != nil {
		t.Fatalf("Could not send message: %s", err)
	}
	beginCacheSizeSuccess := mtest.ScrapeMetricAsInt(t, metrics.ListenAddr, cacheSizeMetricName, cache.Success, 0)
	beginCacheHitSuccess := mtest.ScrapeMetricAsInt(t, metrics.ListenAddr, cacheHitMetricName, cache.Success, 0)
	m = new(dns.Msg)
	m.SetQuestion("www.example.net.", dns.TypeA)
	if _, err = dns.Exchange(m, udp); err != nil {
		t.Fatalf("Could not send message: %s", err)
	}
	got := mtest.ScrapeMetricAsInt(t, metrics.ListenAddr, cacheSizeMetricName, cache.Success, 0)
	if got-beginCacheSizeSuccess != 1 {
		t.Errorf("Expected value %d for %s, but got %d", 1, cacheSizeMetricName, got-beginCacheSizeSuccess)
	}
	if _, err = dns.Exchange(m, udp); err != nil {
		t.Fatalf("Could not send message: %s", err)
	}
	if _, err = dns.Exchange(m, udp); err != nil {
		t.Fatalf("Could not send message: %s", err)
	}
	got = mtest.ScrapeMetricAsInt(t, metrics.ListenAddr, cacheHitMetricName, cache.Success, 0)
	if got-beginCacheHitSuccess != 2 {
		t.Errorf("Expected value %d for %s, but got %d", 2, cacheHitMetricName, got-beginCacheHitSuccess)
	}
}
func TestMetricsAuto(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tmpdir, err := ioutil.TempDir(os.TempDir(), "coredns")
	if err != nil {
		t.Fatal(err)
	}
	corefile := `org:0 {
		auto {
			directory ` + tmpdir + ` db\.(.*) {1} 1
		}
		prometheus localhost:0
	}
`
	i, err := CoreDNSServer(corefile)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}
	udp, _ := CoreDNSServerPorts(i, 0)
	if udp == "" {
		t.Fatalf("Could not get UDP listening port")
	}
	defer i.Stop()
	if err = ioutil.WriteFile(filepath.Join(tmpdir, "db.example.org"), []byte(zoneContent), 0644); err != nil {
		t.Fatal(err)
	}
	time.Sleep(1100 * time.Millisecond)
	m := new(dns.Msg)
	m.SetQuestion("www.example.org.", dns.TypeA)
	if _, err := dns.Exchange(m, udp); err != nil {
		t.Fatalf("Could not send message: %s", err)
	}
	metricName := "coredns_dns_request_count_total"
	data := mtest.Scrape(t, "http://"+metrics.ListenAddr+"/metrics")
	got, _ := mtest.MetricValueLabel(metricName, "example.org.", data)
	if got != "1" {
		t.Errorf("Expected value %s for %s, but got %s", "1", metricName, got)
	}
	os.Remove(filepath.Join(tmpdir, "db.example.org"))
	time.Sleep(1100 * time.Millisecond)
	if _, err := dns.Exchange(m, udp); err != nil {
		t.Fatalf("Could not send message: %s", err)
	}
	data = mtest.Scrape(t, "http://"+metrics.ListenAddr+"/metrics")
	got, _ = mtest.MetricValueLabel(metricName, "example.org.", data)
	if got != "1" {
		t.Errorf("Expected value %s for %s, but got %s", "1", metricName, got)
	}
}
func TestMetricsSeveralBlocs(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cacheSizeMetricName := "coredns_cache_size"
	addrMetrics := "localhost:9155"
	corefile := fmt.Sprintf(`
example.org:0 {
	prometheus %s
	forward . 8.8.8.8:53 {
       force_tcp
    }
}
google.com:0 {
	prometheus %s
	forward . 8.8.8.8:53 {
       force_tcp
    }
	cache
}
`, addrMetrics, addrMetrics)
	i, udp, _, err := CoreDNSServerAndPorts(corefile)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}
	defer i.Stop()
	m := new(dns.Msg)
	m.SetQuestion("google.com.", dns.TypeA)
	if _, err = dns.Exchange(m, udp); err != nil {
		t.Fatalf("Could not send message: %s", err)
	}
	beginCacheSize := mtest.ScrapeMetricAsInt(t, addrMetrics, cacheSizeMetricName, "", 0)
	m = new(dns.Msg)
	m.SetQuestion("www.google.com.", dns.TypeA)
	if _, err = dns.Exchange(m, udp); err != nil {
		t.Fatalf("Could not send message: %s", err)
	}
	endCacheSize := mtest.ScrapeMetricAsInt(t, addrMetrics, cacheSizeMetricName, "", 0)
	if err != nil {
		t.Errorf("Unexpected metric data retrieved for %s : %s", cacheSizeMetricName, err)
	}
	if endCacheSize-beginCacheSize != 1 {
		t.Errorf("Expected metric data retrieved for %s, expected %d, got %d", cacheSizeMetricName, 1, endCacheSize-beginCacheSize)
	}
}
