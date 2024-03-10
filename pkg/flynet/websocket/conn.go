package websocket_net

import (
	"bytes"
	"github.com/gobwas/ws"
	"github.com/linkyaa/fly-im/pkg/flynet/base"
	"github.com/linkyaa/fly-im/pkg/flynet/websocket/wsbase"
	"github.com/linkyaa/fly-im/pkg/pool"
	"github.com/panjf2000/gnet/v2"
	"net"
	"time"
)

type (
	connStatus uint8

	wsConn struct {
		status         connStatus                   //连接升级状态
		connectTime    int64                        //连接到服务器时间,单位s
		conn           gnet.Conn                    //底层连接
		buf            *bytes.Buffer                //用于升级的buf,用完就可以回收了
		bufPool        pool.Pooler[*bytes.Buffer]   //
		cacheFrame     *wsbase.WsFrame              //未完整的frame
		framePool      pool.Pooler[*wsbase.WsFrame] //frame pool
		completeFrames []*wsbase.WsFrame            //engine使用完frames后,重制frames
	}
)

const (
	upgrading connStatus = iota + 1 //升级中
	upgraded                        //升级完成
	auth                            //是否完成认证
	close
)

var _ base.Conn = (*wsConn)(nil)
var emptyHeader ws.Header

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

func (w *wsConn) GetWsFrames() []*wsbase.WsFrame {
	return w.completeFrames
}

func (w *wsConn) ReleaseFrames() {
	w.freeFrames()
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
	if w.status == close {
		return nil
	}

	if w.cacheFrame != nil {
		w.freeFrame(w.cacheFrame)
	}

	if w.buf != nil {
		w.freeBuf()
	}
	w.freeFrames()
	w.status = close
	return nil
}

func (w *wsConn) bufGrow(size int) {
	if w.buf == nil {
		w.buf = w.bufPool.Get()
	}
	w.buf.Grow(size)
}

func (w *wsConn) freeBuf() {
	buf := w.buf
	w.buf = nil
	buf.Reset()
	w.bufPool.Put(buf)
}

func (w *wsConn) getFrame() *wsbase.WsFrame {
	return w.framePool.Get()
}

func (w *wsConn) freeFrame(frame *wsbase.WsFrame) {
	frame.Header = emptyHeader
	w.framePool.Put(frame)
}

func newWsConn(c gnet.Conn, bufPool pool.Pooler[*bytes.Buffer], framePool pool.Pooler[*wsbase.WsFrame]) *wsConn {
	res := &wsConn{
		status:         upgrading,
		connectTime:    time.Now().UTC().Unix(),
		conn:           c,
		bufPool:        bufPool,
		framePool:      framePool,
		completeFrames: make([]*wsbase.WsFrame, 0),
	}
	return res
}
