package rewrite

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

func ipToWire(family int, ipAddr string) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch family {
	case 1:
		return net.ParseIP(ipAddr).To4(), nil
	case 2:
		return net.ParseIP(ipAddr).To16(), nil
	}
	return nil, fmt.Errorf("invalid IP address family (i.e. version) %d", family)
}
func uint16ToWire(data uint16) []byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(data))
	return buf
}
func portToWire(portStr string) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return nil, err
	}
	return uint16ToWire(uint16(port)), nil
}
