package dnsutil

import (
	"net"
	"strings"
)

func ExtractAddressFromReverse(reverseName string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	search := ""
	f := reverse
	switch {
	case strings.HasSuffix(reverseName, IP4arpa):
		search = strings.TrimSuffix(reverseName, IP4arpa)
	case strings.HasSuffix(reverseName, IP6arpa):
		search = strings.TrimSuffix(reverseName, IP6arpa)
		f = reverse6
	default:
		return ""
	}
	return f(strings.Split(search, "."))
}
func IsReverse(name string) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if strings.HasSuffix(name, IP4arpa) {
		return 1
	}
	if strings.HasSuffix(name, IP6arpa) {
		return 2
	}
	return 0
}
func reverse(slice []string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for i := 0; i < len(slice)/2; i++ {
		j := len(slice) - i - 1
		slice[i], slice[j] = slice[j], slice[i]
	}
	ip := net.ParseIP(strings.Join(slice, ".")).To4()
	if ip == nil {
		return ""
	}
	return ip.String()
}
func reverse6(slice []string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for i := 0; i < len(slice)/2; i++ {
		j := len(slice) - i - 1
		slice[i], slice[j] = slice[j], slice[i]
	}
	slice6 := []string{}
	for i := 0; i < len(slice)/4; i++ {
		slice6 = append(slice6, strings.Join(slice[i*4:i*4+4], ""))
	}
	ip := net.ParseIP(strings.Join(slice6, ":")).To16()
	if ip == nil {
		return ""
	}
	return ip.String()
}

const (
	IP4arpa	= ".in-addr.arpa."
	IP6arpa	= ".ip6.arpa."
)
