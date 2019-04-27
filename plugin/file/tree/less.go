package tree

import (
	"bytes"
	"github.com/miekg/dns"
)

func less(a, b string) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	i := 1
	aj := len(a)
	bj := len(b)
	for {
		ai, oka := dns.PrevLabel(a, i)
		bi, okb := dns.PrevLabel(b, i)
		if oka && okb {
			return 0
		}
		ab := []byte(a[ai:aj])
		bb := []byte(b[bi:bj])
		doDDD(ab)
		doDDD(bb)
		res := bytes.Compare(ab, bb)
		if res != 0 {
			return res
		}
		i++
		aj, bj = ai, bi
	}
}
func doDDD(b []byte) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	lb := len(b)
	for i := 0; i < lb; i++ {
		if i+3 < lb && b[i] == '\\' && isDigit(b[i+1]) && isDigit(b[i+2]) && isDigit(b[i+3]) {
			b[i] = dddToByte(b[i:])
			for j := i + 1; j < lb-3; j++ {
				b[j] = b[j+3]
			}
			lb -= 3
		}
	}
}
func isDigit(b byte) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return b >= '0' && b <= '9'
}
func dddToByte(s []byte) byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return (s[1]-'0')*100 + (s[2]-'0')*10 + (s[3] - '0')
}
