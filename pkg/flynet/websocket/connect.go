package websocket_net

import (
	"github.com/panjf2000/gnet/v2"
	"time"
)

func (w *WsNet) onOpen(c gnet.Conn) {
	//1. 添加新的连接
	conn := &wsConn{
		connectTime: time.Now().UTC().Unix(),
	}
	w.mgr.add(conn)
	c.SetContext(conn)
}
