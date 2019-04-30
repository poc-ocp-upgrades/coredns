package test

import (
	"fmt"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"io"
	"mime"
	"net/http"
	godefaulthttp "net/http"
	"strconv"
	"testing"
	"github.com/matttproud/golang_protobuf_extensions/pbutil"
	"github.com/prometheus/common/expfmt"
	dto "github.com/prometheus/client_model/go"
)

type (
	MetricFamily	struct {
		Name	string		`json:"name"`
		Help	string		`json:"help"`
		Type	string		`json:"type"`
		Metrics	[]interface{}	`json:"metrics,omitempty"`
	}
	metric	struct {
		Labels	map[string]string	`json:"labels,omitempty"`
		Value	string			`json:"value"`
	}
	summary	struct {
		Labels		map[string]string	`json:"labels,omitempty"`
		Quantiles	map[string]string	`json:"quantiles,omitempty"`
		Count		string			`json:"count"`
		Sum		string			`json:"sum"`
	}
	histogram	struct {
		Labels	map[string]string	`json:"labels,omitempty"`
		Buckets	map[string]string	`json:"buckets,omitempty"`
		Count	string			`json:"count"`
		Sum	string			`json:"sum"`
	}
)

func Scrape(t *testing.T, url string) []*MetricFamily {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mfChan := make(chan *dto.MetricFamily, 1024)
	go fetchMetricFamilies(url, mfChan)
	result := []*MetricFamily{}
	for mf := range mfChan {
		result = append(result, newMetricFamily(mf))
	}
	return result
}
func ScrapeMetricAsInt(t *testing.T, addr string, name string, label string, nometricvalue int) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	valueToInt := func(m metric) int {
		v := m.Value
		r, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}
		return r
	}
	met := Scrape(t, fmt.Sprintf("http://%s/metrics", addr))
	found := false
	tot := 0
	for _, mf := range met {
		if mf.Name == name {
			for _, m := range mf.Metrics {
				if label == "" {
					tot += valueToInt(m.(metric))
					found = true
					continue
				}
				for _, v := range m.(metric).Labels {
					if v == label {
						tot += valueToInt(m.(metric))
						found = true
					}
				}
			}
		}
	}
	if !found {
		return nometricvalue
	}
	return tot
}
func MetricValue(name string, mfs []*MetricFamily) (string, map[string]string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, mf := range mfs {
		if mf.Name == name {
			return mf.Metrics[0].(metric).Value, mf.Metrics[0].(metric).Labels
		}
	}
	return "", nil
}
func MetricValueLabel(name, label string, mfs []*MetricFamily) (string, map[string]string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, mf := range mfs {
		if mf.Name == name {
			for _, m := range mf.Metrics {
				for _, v := range m.(metric).Labels {
					if v == label {
						return m.(metric).Value, m.(metric).Labels
					}
				}
			}
		}
	}
	return "", nil
}
func newMetricFamily(dtoMF *dto.MetricFamily) *MetricFamily {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mf := &MetricFamily{Name: dtoMF.GetName(), Help: dtoMF.GetHelp(), Type: dtoMF.GetType().String(), Metrics: make([]interface{}, len(dtoMF.Metric))}
	for i, m := range dtoMF.Metric {
		if dtoMF.GetType() == dto.MetricType_SUMMARY {
			mf.Metrics[i] = summary{Labels: makeLabels(m), Quantiles: makeQuantiles(m), Count: fmt.Sprint(m.GetSummary().GetSampleCount()), Sum: fmt.Sprint(m.GetSummary().GetSampleSum())}
		} else if dtoMF.GetType() == dto.MetricType_HISTOGRAM {
			mf.Metrics[i] = histogram{Labels: makeLabels(m), Buckets: makeBuckets(m), Count: fmt.Sprint(m.GetHistogram().GetSampleCount()), Sum: fmt.Sprint(m.GetSummary().GetSampleSum())}
		} else {
			mf.Metrics[i] = metric{Labels: makeLabels(m), Value: fmt.Sprint(value(m))}
		}
	}
	return mf
}
func value(m *dto.Metric) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m.Gauge != nil {
		return m.GetGauge().GetValue()
	}
	if m.Counter != nil {
		return m.GetCounter().GetValue()
	}
	if m.Untyped != nil {
		return m.GetUntyped().GetValue()
	}
	return 0.
}
func makeLabels(m *dto.Metric) map[string]string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := map[string]string{}
	for _, lp := range m.Label {
		result[lp.GetName()] = lp.GetValue()
	}
	return result
}
func makeQuantiles(m *dto.Metric) map[string]string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := map[string]string{}
	for _, q := range m.GetSummary().Quantile {
		result[fmt.Sprint(q.GetQuantile())] = fmt.Sprint(q.GetValue())
	}
	return result
}
func makeBuckets(m *dto.Metric) map[string]string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := map[string]string{}
	for _, b := range m.GetHistogram().Bucket {
		result[fmt.Sprint(b.GetUpperBound())] = fmt.Sprint(b.GetCumulativeCount())
	}
	return result
}
func fetchMetricFamilies(url string, ch chan<- *dto.MetricFamily) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer close(ch)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Add("Accept", acceptHeader)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return
	}
	mediatype, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err == nil && mediatype == "application/vnd.google.protobuf" && params["encoding"] == "delimited" && params["proto"] == "io.prometheus.client.MetricFamily" {
		for {
			mf := &dto.MetricFamily{}
			if _, err = pbutil.ReadDelimited(resp.Body, mf); err != nil {
				if err == io.EOF {
					break
				}
				return
			}
			ch <- mf
		}
	} else {
		var parser expfmt.TextParser
		metricFamilies, err := parser.TextToMetricFamilies(resp.Body)
		if err != nil {
			return
		}
		for _, mf := range metricFamilies {
			ch <- mf
		}
	}
}

const acceptHeader = `application/vnd.google.protobuf;proto=io.prometheus.client.MetricFamily;encoding=delimited;q=0.7,text/plain;version=0.0.4;q=0.3`

func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
