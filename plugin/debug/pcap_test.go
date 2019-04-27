package debug

import (
	"bytes"
	"fmt"
	golog "log"
	"strings"
	"testing"
	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"
)

func msg() *dns.Msg {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := new(dns.Msg)
	m.SetQuestion("example.local.", dns.TypeA)
	m.SetEdns0(4096, true)
	m.Id = 10
	return m
}
func ExampleLogHexdump() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	buf, _ := msg().Pack()
	h := hexdump(buf)
	fmt.Println(string(h))
}
func TestHexdump(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var f bytes.Buffer
	golog.SetOutput(&f)
	log.D = true
	str := "Hi There!"
	Hexdump(msg(), str)
	logged := f.String()
	if !strings.Contains(logged, "[DEBUG] "+str) {
		t.Errorf("The string %s, is not contained in the logged output: %s", str, logged)
	}
}
func TestHexdumpf(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var f bytes.Buffer
	golog.SetOutput(&f)
	log.D = true
	str := "Hi There!"
	Hexdumpf(msg(), "%s %d", str, 10)
	logged := f.String()
	if !strings.Contains(logged, "[DEBUG] "+fmt.Sprintf("%s %d", str, 10)) {
		t.Errorf("The string %s %d, is not contained in the logged output: %s", str, 10, logged)
	}
}
func TestNoDebug(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var f bytes.Buffer
	golog.SetOutput(&f)
	log.D = false
	str := "Hi There!"
	Hexdumpf(msg(), "%s %d", str, 10)
	if len(f.Bytes()) != 0 {
		t.Errorf("Expected no output, got %d bytes", len(f.Bytes()))
	}
}
