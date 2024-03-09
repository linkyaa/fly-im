package websocket_net

import (
	"context"
	"github.com/linkyaa/fly-im/pkg/flynet/base"
	"github.com/linkyaa/fly-im/pkg/logx"
	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
	"sync"
	"time"
)

type (
	WsNet struct {
		gnet.Engine
		wsAddr  string
		opt     *base.Options
		handler base.EventHandler
		lock    sync.RWMutex
		mgr     *connMgr //连接管理
	}
)

func (w *WsNet) OnBoot(eng gnet.Engine) (action gnet.Action) {
	w.Engine = eng
	return gnet.None
}

func (w *WsNet) OnShutdown(_ gnet.Engine) {
	logx.Logger.Warn("websocket server shutdown")
}

func (w *WsNet) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	w.onOpen(c)
	return nil, gnet.None
}

func (w *WsNet) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	conn, ok := c.Context().(*wsConn)
	if !ok {
		return gnet.Close
	}
	w.handler.OnClose(conn, err)
	return gnet.Close
}

func (w *WsNet) OnTraffic(c gnet.Conn) (action gnet.Action) {
	//TODO:升级连接
	return gnet.None
}

func (w *WsNet) OnTick() (delay time.Duration, action gnet.Action) {
	return time.Minute, gnet.None
}

func (w *WsNet) Run() {
	err := gnet.Run(w, w.wsAddr,
		gnet.WithTicker(true),
		gnet.WithMulticore(true),
	)
	if err != nil {
		logx.Logger.Error("gnet run err", zap.Error(err))
	}
}

func (w *WsNet) Stop() {
	err := w.Engine.Stop(context.Background())
	if err != nil {
		logx.Logger.Error("stop gnet base err", zap.Error(err))
	}
}

func New(addr string, handler base.EventHandler, opt *base.Options) *WsNet {
	res := &WsNet{
		wsAddr:  addr,
		opt:     opt,
		handler: handler,
		mgr:     newConnMgr(),
	}
	return res
}
