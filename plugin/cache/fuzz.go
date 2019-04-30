package cache

import (
	"github.com/coredns/coredns/plugin/pkg/fuzz"
)

func Fuzz(data []byte) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fuzz.Do(New(), data)
}
