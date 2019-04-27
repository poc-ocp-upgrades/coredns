package chaos

import (
	"github.com/coredns/coredns/plugin/pkg/fuzz"
)

func Fuzz(data []byte) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := Chaos{}
	return fuzz.Do(c, data)
}
