package rewrite

import (
	"context"
	"fmt"
	"strings"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Result int

const (
	RewriteIgnored	Result	= iota
	RewriteDone
)
const (
	Stop		= "stop"
	Continue	= "continue"
)

type Rewrite struct {
	Next		plugin.Handler
	Rules		[]Rule
	noRevert	bool
}

func (rw Rewrite) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	wr := NewResponseReverter(w, r)
	state := request.Request{W: w, Req: r}
	for _, rule := range rw.Rules {
		switch result := rule.Rewrite(ctx, state); result {
		case RewriteDone:
			if !validName(state.Req.Question[0].Name) {
				x := state.Req.Question[0].Name
				log.Errorf("Invalid name after rewrite: %s", x)
				state.Req.Question[0] = wr.originalQuestion
				return dns.RcodeServerFailure, fmt.Errorf("invalid name after rewrite: %s", x)
			}
			respRule := rule.GetResponseRule()
			if respRule.Active == true {
				wr.ResponseRewrite = true
				wr.ResponseRules = append(wr.ResponseRules, respRule)
			}
			if rule.Mode() == Stop {
				if rw.noRevert {
					return plugin.NextOrFailure(rw.Name(), rw.Next, ctx, w, r)
				}
				return plugin.NextOrFailure(rw.Name(), rw.Next, ctx, wr, r)
			}
		case RewriteIgnored:
			break
		}
	}
	if rw.noRevert || len(wr.ResponseRules) == 0 {
		return plugin.NextOrFailure(rw.Name(), rw.Next, ctx, w, r)
	}
	return plugin.NextOrFailure(rw.Name(), rw.Next, ctx, wr, r)
}
func (rw Rewrite) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "rewrite"
}

type Rule interface {
	Rewrite(ctx context.Context, state request.Request) Result
	Mode() string
	GetResponseRule() ResponseRule
}

func newRule(args ...string) (Rule, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(args) == 0 {
		return nil, fmt.Errorf("no rule type specified for rewrite")
	}
	arg0 := strings.ToLower(args[0])
	var ruleType string
	var expectNumArgs, startArg int
	mode := Stop
	switch arg0 {
	case Continue:
		mode = Continue
		ruleType = strings.ToLower(args[1])
		expectNumArgs = len(args) - 1
		startArg = 2
	case Stop:
		ruleType = strings.ToLower(args[1])
		expectNumArgs = len(args) - 1
		startArg = 2
	default:
		ruleType = arg0
		expectNumArgs = len(args)
		startArg = 1
	}
	switch ruleType {
	case "answer":
		return nil, fmt.Errorf("response rewrites must begin with a name rule")
	case "name":
		return newNameRule(mode, args[startArg:]...)
	case "class":
		if expectNumArgs != 3 {
			return nil, fmt.Errorf("%s rules must have exactly two arguments", ruleType)
		}
		return newClassRule(mode, args[startArg:]...)
	case "type":
		if expectNumArgs != 3 {
			return nil, fmt.Errorf("%s rules must have exactly two arguments", ruleType)
		}
		return newTypeRule(mode, args[startArg:]...)
	case "edns0":
		return newEdns0Rule(mode, args[startArg:]...)
	case "ttl":
		return newTtlRule(mode, args[startArg:]...)
	default:
		return nil, fmt.Errorf("invalid rule type %q", args[0])
	}
}
