package rewrite

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"github.com/coredns/coredns/plugin/pkg/replacer"
	"github.com/miekg/dns"
)

const (
	Is		= "is"
	Not		= "not"
	Has		= "has"
	NotHas		= "not_has"
	StartsWith	= "starts_with"
	EndsWith	= "ends_with"
	Match		= "match"
	NotMatch	= "not_match"
)

func newReplacer(r *dns.Msg) replacer.Replacer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return replacer.New(context.TODO(), r, nil, "")
}

type condition func(string, string) bool

var conditions = map[string]condition{Is: isFunc, Not: notFunc, Has: hasFunc, NotHas: notHasFunc, StartsWith: startsWithFunc, EndsWith: endsWithFunc, Match: matchFunc, NotMatch: notMatchFunc}

func isFunc(a, b string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return a == b
}
func notFunc(a, b string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return a != b
}
func hasFunc(a, b string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return strings.Contains(a, b)
}
func notHasFunc(a, b string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return !strings.Contains(a, b)
}
func startsWithFunc(a, b string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return strings.HasPrefix(a, b)
}
func endsWithFunc(a, b string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return strings.HasSuffix(a, b)
}
func matchFunc(a, b string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	matched, _ := regexp.MatchString(b, a)
	return matched
}
func notMatchFunc(a, b string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	matched, _ := regexp.MatchString(b, a)
	return !matched
}

type If struct {
	A		string
	Operator	string
	B		string
}

func (i If) True(r *dns.Msg) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c, ok := conditions[i.Operator]; ok {
		a, b := i.A, i.B
		if r != nil {
			replacer := newReplacer(r)
			a = replacer.Replace(i.A)
			b = replacer.Replace(i.B)
		}
		return c(a, b)
	}
	return false
}
func NewIf(a, operator, b string) (If, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if _, ok := conditions[operator]; !ok {
		return If{}, fmt.Errorf("invalid operator %v", operator)
	}
	return If{A: a, Operator: operator, B: b}, nil
}
