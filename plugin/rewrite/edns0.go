package rewrite

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"github.com/coredns/coredns/plugin/metadata"
	"github.com/coredns/coredns/plugin/pkg/edns"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type edns0LocalRule struct {
	mode	string
	action	string
	code	uint16
	data	[]byte
}
type edns0VariableRule struct {
	mode		string
	action		string
	code		uint16
	variable	string
}
type edns0NsidRule struct {
	mode	string
	action	string
}

func setupEdns0Opt(r *dns.Msg) *dns.OPT {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	o := r.IsEdns0()
	if o == nil {
		r.SetEdns0(4096, false)
		o = r.IsEdns0()
	}
	return o
}
func (rule *edns0NsidRule) Rewrite(ctx context.Context, state request.Request) Result {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	o := setupEdns0Opt(state.Req)
	for _, s := range o.Option {
		if e, ok := s.(*dns.EDNS0_NSID); ok {
			if rule.action == Replace || rule.action == Set {
				e.Nsid = ""
				return RewriteDone
			}
		}
	}
	if rule.action == Append || rule.action == Set {
		o.Option = append(o.Option, &dns.EDNS0_NSID{Code: dns.EDNS0NSID, Nsid: ""})
		return RewriteDone
	}
	return RewriteIgnored
}
func (rule *edns0NsidRule) Mode() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return rule.mode
}
func (rule *edns0NsidRule) GetResponseRule() ResponseRule {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ResponseRule{}
}
func (rule *edns0LocalRule) Rewrite(ctx context.Context, state request.Request) Result {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	o := setupEdns0Opt(state.Req)
	for _, s := range o.Option {
		if e, ok := s.(*dns.EDNS0_LOCAL); ok {
			if rule.code == e.Code {
				if rule.action == Replace || rule.action == Set {
					e.Data = rule.data
					return RewriteDone
				}
			}
		}
	}
	if rule.action == Append || rule.action == Set {
		o.Option = append(o.Option, &dns.EDNS0_LOCAL{Code: rule.code, Data: rule.data})
		return RewriteDone
	}
	return RewriteIgnored
}
func (rule *edns0LocalRule) Mode() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return rule.mode
}
func (rule *edns0LocalRule) GetResponseRule() ResponseRule {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ResponseRule{}
}
func newEdns0Rule(mode string, args ...string) (Rule, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(args) < 2 {
		return nil, fmt.Errorf("too few arguments for an EDNS0 rule")
	}
	ruleType := strings.ToLower(args[0])
	action := strings.ToLower(args[1])
	switch action {
	case Append:
	case Replace:
	case Set:
	default:
		return nil, fmt.Errorf("invalid action: %q", action)
	}
	switch ruleType {
	case "local":
		if len(args) != 4 {
			return nil, fmt.Errorf("EDNS0 local rules require exactly three args")
		}
		if strings.HasPrefix(args[3], "{") && strings.HasSuffix(args[3], "}") {
			return newEdns0VariableRule(mode, action, args[2], args[3])
		}
		return newEdns0LocalRule(mode, action, args[2], args[3])
	case "nsid":
		if len(args) != 2 {
			return nil, fmt.Errorf("EDNS0 NSID rules do not accept args")
		}
		return &edns0NsidRule{mode: mode, action: action}, nil
	case "subnet":
		if len(args) != 4 {
			return nil, fmt.Errorf("EDNS0 subnet rules require exactly three args")
		}
		return newEdns0SubnetRule(mode, action, args[2], args[3])
	default:
		return nil, fmt.Errorf("invalid rule type %q", ruleType)
	}
}
func newEdns0LocalRule(mode, action, code, data string) (*edns0LocalRule, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	c, err := strconv.ParseUint(code, 0, 16)
	if err != nil {
		return nil, err
	}
	decoded := []byte(data)
	if strings.HasPrefix(data, "0x") {
		decoded, err = hex.DecodeString(data[2:])
		if err != nil {
			return nil, err
		}
	}
	edns.SetSupportedOption(uint16(c))
	return &edns0LocalRule{mode: mode, action: action, code: uint16(c), data: decoded}, nil
}
func newEdns0VariableRule(mode, action, code, variable string) (*edns0VariableRule, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	c, err := strconv.ParseUint(code, 0, 16)
	if err != nil {
		return nil, err
	}
	if !isValidVariable(variable) {
		return nil, fmt.Errorf("unsupported variable name %q", variable)
	}
	edns.SetSupportedOption(uint16(c))
	return &edns0VariableRule{mode: mode, action: action, code: uint16(c), variable: variable}, nil
}
func (rule *edns0VariableRule) ruleData(ctx context.Context, state request.Request) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch rule.variable {
	case queryName:
		return []byte(state.QName()), nil
	case queryType:
		return uint16ToWire(state.QType()), nil
	case clientIP:
		return ipToWire(state.Family(), state.IP())
	case serverIP:
		return ipToWire(state.Family(), state.LocalIP())
	case clientPort:
		return portToWire(state.Port())
	case serverPort:
		return portToWire(state.LocalPort())
	case protocol:
		return []byte(state.Proto()), nil
	}
	fetcher := metadata.ValueFunc(ctx, rule.variable[1:len(rule.variable)-1])
	if fetcher != nil {
		value := fetcher()
		if len(value) > 0 {
			return []byte(value), nil
		}
	}
	return nil, fmt.Errorf("unable to extract data for variable %s", rule.variable)
}
func (rule *edns0VariableRule) Rewrite(ctx context.Context, state request.Request) Result {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	data, err := rule.ruleData(ctx, state)
	if err != nil || data == nil {
		return RewriteIgnored
	}
	o := setupEdns0Opt(state.Req)
	for _, s := range o.Option {
		if e, ok := s.(*dns.EDNS0_LOCAL); ok {
			if rule.code == e.Code {
				if rule.action == Replace || rule.action == Set {
					e.Data = data
					return RewriteDone
				}
				return RewriteIgnored
			}
		}
	}
	if rule.action == Append || rule.action == Set {
		o.Option = append(o.Option, &dns.EDNS0_LOCAL{Code: rule.code, Data: data})
		return RewriteDone
	}
	return RewriteIgnored
}
func (rule *edns0VariableRule) Mode() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return rule.mode
}
func (rule *edns0VariableRule) GetResponseRule() ResponseRule {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ResponseRule{}
}
func isValidVariable(variable string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch variable {
	case queryName, queryType, clientIP, clientPort, protocol, serverIP, serverPort:
		return true
	}
	if strings.HasPrefix(variable, "{") && strings.HasSuffix(variable, "}") && metadata.IsLabel(variable[1:len(variable)-1]) {
		return true
	}
	return false
}

type edns0SubnetRule struct {
	mode		string
	v4BitMaskLen	uint8
	v6BitMaskLen	uint8
	action		string
}

func newEdns0SubnetRule(mode, action, v4BitMaskLen, v6BitMaskLen string) (*edns0SubnetRule, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v4Len, err := strconv.ParseUint(v4BitMaskLen, 0, 16)
	if err != nil {
		return nil, err
	}
	if v4Len > net.IPv4len*8 {
		return nil, fmt.Errorf("invalid IPv4 bit mask length %d", v4Len)
	}
	v6Len, err := strconv.ParseUint(v6BitMaskLen, 0, 16)
	if err != nil {
		return nil, err
	}
	if v6Len > net.IPv6len*8 {
		return nil, fmt.Errorf("invalid IPv6 bit mask length %d", v6Len)
	}
	return &edns0SubnetRule{mode: mode, action: action, v4BitMaskLen: uint8(v4Len), v6BitMaskLen: uint8(v6Len)}, nil
}
func (rule *edns0SubnetRule) fillEcsData(state request.Request, ecs *dns.EDNS0_SUBNET) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	family := state.Family()
	if (family != 1) && (family != 2) {
		return fmt.Errorf("unable to fill data for EDNS0 subnet due to invalid IP family")
	}
	ecs.Family = uint16(family)
	ecs.SourceScope = 0
	ipAddr := state.IP()
	switch family {
	case 1:
		ipv4Mask := net.CIDRMask(int(rule.v4BitMaskLen), 32)
		ipv4Addr := net.ParseIP(ipAddr)
		ecs.SourceNetmask = rule.v4BitMaskLen
		ecs.Address = ipv4Addr.Mask(ipv4Mask).To4()
	case 2:
		ipv6Mask := net.CIDRMask(int(rule.v6BitMaskLen), 128)
		ipv6Addr := net.ParseIP(ipAddr)
		ecs.SourceNetmask = rule.v6BitMaskLen
		ecs.Address = ipv6Addr.Mask(ipv6Mask).To16()
	}
	return nil
}
func (rule *edns0SubnetRule) Rewrite(ctx context.Context, state request.Request) Result {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	o := setupEdns0Opt(state.Req)
	for _, s := range o.Option {
		if e, ok := s.(*dns.EDNS0_SUBNET); ok {
			if rule.action == Replace || rule.action == Set {
				if rule.fillEcsData(state, e) == nil {
					return RewriteDone
				}
			}
			return RewriteIgnored
		}
	}
	if rule.action == Append || rule.action == Set {
		opt := &dns.EDNS0_SUBNET{Code: dns.EDNS0SUBNET}
		if rule.fillEcsData(state, opt) == nil {
			o.Option = append(o.Option, opt)
			return RewriteDone
		}
	}
	return RewriteIgnored
}
func (rule *edns0SubnetRule) Mode() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return rule.mode
}
func (rule *edns0SubnetRule) GetResponseRule() ResponseRule {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ResponseRule{}
}

const (
	Replace	= "replace"
	Set	= "set"
	Append	= "append"
)
const (
	queryName	= "{qname}"
	queryType	= "{qtype}"
	clientIP	= "{client_ip}"
	clientPort	= "{client_port}"
	protocol	= "{protocol}"
	serverIP	= "{server_ip}"
	serverPort	= "{server_port}"
)
