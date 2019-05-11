package dnsserver

import "net"

func listen(network, addr string) (net.Listener, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return net.Listen(network, addr)
}
func listenPacket(network, addr string) (net.PacketConn, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return net.ListenPacket(network, addr)
}
