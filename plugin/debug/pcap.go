package debug

import (
	"bytes"
	"fmt"
	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"
)

func Hexdump(m *dns.Msg, v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !log.D {
		return
	}
	buf, _ := m.Pack()
	if len(buf) == 0 {
		return
	}
	out := "\n" + string(hexdump(buf))
	v = append(v, out)
	log.Debug(v...)
}
func Hexdumpf(m *dns.Msg, format string, v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !log.D {
		return
	}
	buf, _ := m.Pack()
	if len(buf) == 0 {
		return
	}
	format += "\n%s"
	v = append(v, hexdump(buf))
	log.Debugf(format, v...)
}
func hexdump(data []byte) []byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	b := new(bytes.Buffer)
	newline := ""
	for i := 0; i < len(data); i++ {
		if i%16 == 0 {
			fmt.Fprintf(b, "%s%s%06x", newline, prefix, i)
			newline = "\n"
		}
		fmt.Fprintf(b, " %02x", data[i])
	}
	fmt.Fprintf(b, "\n%s%06x", prefix, len(data))
	return b.Bytes()
}

const prefix = "debug: "
