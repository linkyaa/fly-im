package websocket_net

import (
	"github.com/linkyaa/fly-im/pkg/flynet/websocket/wsconn"
	"github.com/linkyaa/fly-im/pkg/logx"
	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
	"time"
)

func (w *WsNet) OnBoot(eng gnet.Engine) (action gnet.Action) {
	w.Engine = eng
	return gnet.None
}

func (w *WsNet) OnShutdown(_ gnet.Engine) {
	logx.Warn("ws onShutdown")
}

func (w *WsNet) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	conn := wsconn.NewWsConn(c, w.bufPool, w.framePool)
	c.SetContext(conn)
	w.handler.OnConnect(conn)
	return nil, gnet.None
}

func (w *WsNet) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	conn, ok := c.Context().(*wsconn.WsConn)
	if !ok {
		logx.Warn("ws onClose 无法获取context")
		return gnet.Close
	}

	w.handler.OnClose(conn, err)
	return gnet.Close
}

func (w *WsNet) OnTraffic(c gnet.Conn) (action gnet.Action) {
	conn, ok := c.Context().(*wsconn.WsConn)
	if !ok {
		logx.Warn("ws onTraffic无法获取context")
		return gnet.Close
	}

	//1. 判断是否升级
	done, err := conn.TryUpgrade()
	if err != nil {
		logx.Error("ws tryUpgrade err", zap.Error(err))
		return gnet.Close
	}

	if !done {
		return gnet.None
	}

	//2. 解码消息
	ok, err = conn.DecodeFrame()
	if err != nil {
		logx.Debug("ws decodeFrame err", zap.Error(err))
		return gnet.Close
	}

	if ok {
		w.handler.OnData(conn)
	}
	return gnet.None
}

func (w *WsNet) OnTick() (delay time.Duration, action gnet.Action) {
	return time.Minute, gnet.None
}
