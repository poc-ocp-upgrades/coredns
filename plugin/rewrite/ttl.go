package rewrite

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
)

type exactTtlRule struct {
	NextAction	string
	From		string
	ResponseRule
}
type prefixTtlRule struct {
	NextAction	string
	Prefix		string
	ResponseRule
}
type suffixTtlRule struct {
	NextAction	string
	Suffix		string
	ResponseRule
}
type substringTtlRule struct {
	NextAction	string
	Substring	string
	ResponseRule
}
type regexTtlRule struct {
	NextAction	string
	Pattern		*regexp.Regexp
	ResponseRule
}

func (rule *exactTtlRule) Rewrite(ctx context.Context, state request.Request) Result {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if rule.From == state.Name() {
		return RewriteDone
	}
	return RewriteIgnored
}
func (rule *prefixTtlRule) Rewrite(ctx context.Context, state request.Request) Result {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if strings.HasPrefix(state.Name(), rule.Prefix) {
		return RewriteDone
	}
	return RewriteIgnored
}
func (rule *suffixTtlRule) Rewrite(ctx context.Context, state request.Request) Result {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if strings.HasSuffix(state.Name(), rule.Suffix) {
		return RewriteDone
	}
	return RewriteIgnored
}
func (rule *substringTtlRule) Rewrite(ctx context.Context, state request.Request) Result {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if strings.Contains(state.Name(), rule.Substring) {
		return RewriteDone
	}
	return RewriteIgnored
}
func (rule *regexTtlRule) Rewrite(ctx context.Context, state request.Request) Result {
	_logClusterCodePath()
	defer _logClusterCodePath()
	regexGroups := rule.Pattern.FindStringSubmatch(state.Name())
	if len(regexGroups) == 0 {
		return RewriteIgnored
	}
	return RewriteDone
}
func newTtlRule(nextAction string, args ...string) (Rule, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(args) < 2 {
		return nil, fmt.Errorf("too few (%d) arguments for a ttl rule", len(args))
	}
	var s string
	if len(args) == 2 {
		s = args[1]
	}
	if len(args) == 3 {
		s = args[2]
	}
	ttl, valid := isValidTtl(s)
	if valid == false {
		return nil, fmt.Errorf("invalid TTL '%s' for a ttl rule", s)
	}
	if len(args) == 3 {
		switch strings.ToLower(args[0]) {
		case ExactMatch:
			return &exactTtlRule{nextAction, plugin.Name(args[1]).Normalize(), ResponseRule{Active: true, Type: "ttl", Ttl: ttl}}, nil
		case PrefixMatch:
			return &prefixTtlRule{nextAction, plugin.Name(args[1]).Normalize(), ResponseRule{Active: true, Type: "ttl", Ttl: ttl}}, nil
		case SuffixMatch:
			return &suffixTtlRule{nextAction, plugin.Name(args[1]).Normalize(), ResponseRule{Active: true, Type: "ttl", Ttl: ttl}}, nil
		case SubstringMatch:
			return &substringTtlRule{nextAction, plugin.Name(args[1]).Normalize(), ResponseRule{Active: true, Type: "ttl", Ttl: ttl}}, nil
		case RegexMatch:
			regexPattern, err := regexp.Compile(args[1])
			if err != nil {
				return nil, fmt.Errorf("Invalid regex pattern in a ttl rule: %s", args[1])
			}
			return &regexTtlRule{nextAction, regexPattern, ResponseRule{Active: true, Type: "ttl", Ttl: ttl}}, nil
		default:
			return nil, fmt.Errorf("A ttl rule supports only exact, prefix, suffix, substring, and regex name matching")
		}
	}
	if len(args) > 3 {
		return nil, fmt.Errorf("many few arguments for a ttl rule")
	}
	return &exactTtlRule{nextAction, plugin.Name(args[0]).Normalize(), ResponseRule{Active: true, Type: "ttl", Ttl: ttl}}, nil
}
func (rule *exactTtlRule) Mode() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return rule.NextAction
}
func (rule *prefixTtlRule) Mode() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return rule.NextAction
}
func (rule *suffixTtlRule) Mode() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return rule.NextAction
}
func (rule *substringTtlRule) Mode() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return rule.NextAction
}
func (rule *regexTtlRule) Mode() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return rule.NextAction
}
func (rule *exactTtlRule) GetResponseRule() ResponseRule {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return rule.ResponseRule
}
func (rule *prefixTtlRule) GetResponseRule() ResponseRule {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return rule.ResponseRule
}
func (rule *suffixTtlRule) GetResponseRule() ResponseRule {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return rule.ResponseRule
}
func (rule *substringTtlRule) GetResponseRule() ResponseRule {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return rule.ResponseRule
}
func (rule *regexTtlRule) GetResponseRule() ResponseRule {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return rule.ResponseRule
}
func isValidTtl(v string) (uint32, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	i, err := strconv.Atoi(v)
	if err != nil {
		return uint32(0), false
	}
	if i > 2147483647 {
		return uint32(0), false
	}
	if i < 0 {
		return uint32(0), false
	}
	return uint32(i), true
}
