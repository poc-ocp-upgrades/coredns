package rewrite

import (
	"github.com/coredns/coredns/plugin/pkg/fuzz"
	"github.com/mholt/caddy"
)

func Fuzz(data []byte) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := caddy.NewTestController("dns", "rewrite edns0 subnet set 24 56")
	rules, err := rewriteParse(c)
	if err != nil {
		return 0
	}
	r := Rewrite{Rules: rules}
	return fuzz.Do(r, data)
}
