package dnstapio

import (
	"net"
	"sync/atomic"
	"time"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	tap "github.com/dnstap/golang-dnstap"
	fs "github.com/farsightsec/golang-framestream"
)

var log = clog.NewWithPlugin("dnstap")

const (
	tcpWriteBufSize	= 1024 * 1024
	tcpTimeout	= 4 * time.Second
	flushTimeout	= 1 * time.Second
	queueSize	= 10000
)

type dnstapIO struct {
	endpoint	string
	socket		bool
	conn		net.Conn
	enc		*dnstapEncoder
	queue		chan tap.Dnstap
	dropped		uint32
	quit		chan struct{}
}

func New(endpoint string, socket bool) DnstapIO {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &dnstapIO{endpoint: endpoint, socket: socket, enc: newDnstapEncoder(&fs.EncoderOptions{ContentType: []byte("protobuf:dnstap.Dnstap"), Bidirectional: true}), queue: make(chan tap.Dnstap, queueSize), quit: make(chan struct{})}
}

type DnstapIO interface {
	Connect()
	Dnstap(payload tap.Dnstap)
	Close()
}

func (dio *dnstapIO) newConnect() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var err error
	if dio.socket {
		if dio.conn, err = net.Dial("unix", dio.endpoint); err != nil {
			return err
		}
	} else {
		if dio.conn, err = net.DialTimeout("tcp", dio.endpoint, tcpTimeout); err != nil {
			return err
		}
		if tcpConn, ok := dio.conn.(*net.TCPConn); ok {
			tcpConn.SetWriteBuffer(tcpWriteBufSize)
			tcpConn.SetNoDelay(false)
		}
	}
	return dio.enc.resetWriter(dio.conn)
}
func (dio *dnstapIO) Connect() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := dio.newConnect(); err != nil {
		log.Error("No connection to dnstap endpoint")
	}
	go dio.serve()
}
func (dio *dnstapIO) Dnstap(payload tap.Dnstap) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	select {
	case dio.queue <- payload:
	default:
		atomic.AddUint32(&dio.dropped, 1)
	}
}
func (dio *dnstapIO) closeConnection() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	dio.enc.close()
	if dio.conn != nil {
		dio.conn.Close()
		dio.conn = nil
	}
}
func (dio *dnstapIO) Close() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	close(dio.quit)
}
func (dio *dnstapIO) flushBuffer() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if dio.conn == nil {
		if err := dio.newConnect(); err != nil {
			return
		}
		log.Info("Reconnected to dnstap")
	}
	if err := dio.enc.flushBuffer(); err != nil {
		log.Warningf("Connection lost: %s", err)
		dio.closeConnection()
		if err := dio.newConnect(); err != nil {
			log.Errorf("Cannot connect to dnstap: %s", err)
		} else {
			log.Info("Reconnected to dnstap")
		}
	}
}
func (dio *dnstapIO) write(payload *tap.Dnstap) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := dio.enc.writeMsg(payload); err != nil {
		atomic.AddUint32(&dio.dropped, 1)
	}
}
func (dio *dnstapIO) serve() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	timeout := time.After(flushTimeout)
	for {
		select {
		case <-dio.quit:
			dio.flushBuffer()
			dio.closeConnection()
			return
		case payload := <-dio.queue:
			dio.write(&payload)
		case <-timeout:
			if dropped := atomic.SwapUint32(&dio.dropped, 0); dropped > 0 {
				log.Warningf("Dropped dnstap messages: %d", dropped)
			}
			dio.flushBuffer()
			timeout = time.After(flushTimeout)
		}
	}
}
