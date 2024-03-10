package websocket_net

import (
	"bytes"
	"context"
	"github.com/linkyaa/fly-im/pkg/flynet/base"
	"github.com/linkyaa/fly-im/pkg/flynet/websocket/wsbase"
	"github.com/linkyaa/fly-im/pkg/logx"
	"github.com/linkyaa/fly-im/pkg/pool"
	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
	"io"
	"net/url"
	"time"
)

type (
	WsNet struct {
		gnet.Engine
		wsAddr    string
		opt       *base.Options
		handler   base.EventHandler
		mgr       *connMgr //连接管理
		bufPool   pool.Pooler[*bytes.Buffer]
		framePool pool.Pooler[*wsbase.WsFrame]
	}

	upgrader struct {
		io.Reader
		io.Writer
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

	_ = conn.Close()
	w.handler.OnClose(conn, err)
	return gnet.Close
}

func (w *WsNet) OnTraffic(c gnet.Conn) (action gnet.Action) {
	conn, ok := c.Context().(*wsConn)
	if !ok {
		return gnet.Close
	}

	if conn.status == upgrading {
		ok, err := conn.procUpgrade()
		if err != nil {
			logx.Logger.Error("upgrade err", zap.Error(err))
			return gnet.Close
		}

		if !ok {
			return gnet.None
		}

		w.handler.OnConnect(conn)
		return
	}

	//TODO：检查是否有消息

	err := conn.procWsFrames()
	if err != nil {
		return gnet.Close
	}

	w.handler.OnData(conn)
	return gnet.None
}

func (w *WsNet) OnTick() (delay time.Duration, action gnet.Action) {
	return time.Minute, gnet.None
}

func (w *WsNet) Run() {

	parse, err := url.Parse(w.wsAddr)
	if logx.EnableDebug() {
		logx.Logger.Info("websocket server", zap.String("scheme", parse.Scheme),
			zap.String("host", parse.Host),
			zap.String("path", parse.Path),
			zap.String("port", parse.Port()),
			zap.String("query", parse.RawQuery),
			zap.String("fragment", parse.Fragment),
		)
	}
	//TODO:支持wss协议
	if parse.Scheme == "wss" {
		logx.Logger.Error("暂不支持wss协议")
		return
	}

	err = gnet.Run(w, parse.Host,
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
		logx.Logger.Error("stop gnet err", zap.Error(err))
	}
}

func New(addr string, handler base.EventHandler, opt *base.Options) *WsNet {
	res := &WsNet{
		wsAddr:  addr,
		opt:     opt,
		handler: handler,
		mgr:     newConnMgr(),
		bufPool: pool.NewStdPool[*bytes.Buffer](func() any {
			return &bytes.Buffer{}
		}),
		framePool: pool.NewStdPool[*wsbase.WsFrame](func() any {
			return wsbase.NewWsFrame(4096)
		}),
	}
	return res
}
