package websocket_net

import (
	"github.com/linkyaa/fly-im/pkg/flynet/base"
	"github.com/panjf2000/gnet/v2"
	"net"
)

type (
	wsConn struct {
		conn        gnet.Conn
		connectTime int64 //连接到服务器时间,单位s
	}
)

var c base.Conn = (*wsConn)(nil)

func (w *wsConn) SetAuth(auth bool) {
	//TODO implement me
	panic("implement me")
}

func (w *wsConn) IsAuth() bool {
	//TODO implement me
	panic("implement me")
}

func (w *wsConn) RemoteAddr() net.Addr {
	//TODO implement me
	panic("implement me")
}

func (w *wsConn) LocalAddr() net.Addr {
	//TODO implement me
	panic("implement me")
}

func (w *wsConn) UserId() int64 {
	//TODO implement me
	panic("implement me")
}

func (w *wsConn) GetFrames() ([]*base.Frame, error) {
	//TODO implement me
	panic("implement me")
}

func (w *wsConn) ReleaseFrames() {
	//TODO implement me
	panic("implement me")
}

func (w *wsConn) Write(frame *base.Frame) error {
	//TODO implement me
	panic("implement me")
}

func (w *wsConn) AsyncWrite(frame *base.Frame) error {
	//TODO implement me
	panic("implement me")
}

func (w *wsConn) Close() error {
	//TODO implement me
	panic("implement me")
}
