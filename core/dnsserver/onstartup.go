package dnsserver

import "fmt"

func startUpZones(protocol, addr string, zones map[string]*Config) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := ""
	for zone := range zones {
		_, ip, port, err := SplitProtocolHostPort(addr)
		if err != nil {
			s += fmt.Sprintln(protocol + zone + ":" + addr)
			continue
		}
		if ip == "" {
			s += fmt.Sprintln(protocol + zone + ":" + port)
			continue
		}
		s += fmt.Sprintln(protocol + zone + ":" + port + " on " + ip)
	}
	return s
}
