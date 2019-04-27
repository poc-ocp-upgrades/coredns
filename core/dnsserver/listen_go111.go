package dnsserver

import (
	"context"
	"net"
	"syscall"
	"github.com/coredns/coredns/plugin/pkg/log"
	"golang.org/x/sys/unix"
)

func reuseportControl(network, address string, c syscall.RawConn) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.Control(func(fd uintptr) {
		if err := unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
			log.Warningf("Failed to set SO_REUSEPORT on socket: %s", err)
		}
	})
	return nil
}
func listen(network, addr string) (net.Listener, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	lc := net.ListenConfig{Control: reuseportControl}
	return lc.Listen(context.Background(), network, addr)
}
func listenPacket(network, addr string) (net.PacketConn, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	lc := net.ListenConfig{Control: reuseportControl}
	return lc.ListenPacket(context.Background(), network, addr)
}
