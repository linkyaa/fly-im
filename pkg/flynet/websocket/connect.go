package websocket_net

import (
	"github.com/panjf2000/gnet/v2"
)

func (w *WsNet) onOpen(c gnet.Conn) {
	//1. 添加新的连接
	conn := newWsConn(c, w.bufPool, w.framePool)
	w.mgr.add(conn)
	c.SetContext(conn)
}
