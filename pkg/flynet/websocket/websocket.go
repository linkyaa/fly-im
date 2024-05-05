package websocket_net

import (
	"bytes"
	"context"
	"github.com/linkyaa/fly-im/pkg/flynet"
	"github.com/linkyaa/fly-im/pkg/flynet/frame"
	"github.com/linkyaa/fly-im/pkg/logx"
	"github.com/linkyaa/fly-im/pkg/pool"
	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
	"net/url"
)

type (
	WsNet struct {
		gnet.Engine                             //gnet engine
		wsAddr      string                      //监听地址
		opt         *flynet.Options             //运行参数
		handler     flynet.EventHandler         //事件处理
		bufPool     pool.Pooler[*bytes.Buffer]  //bytes内存池
		framePool   pool.Pooler[*frame.WsFrame] //frame内存池
	}
)

var (
	_ gnet.EventHandler = (*WsNet)(nil)
)

func (w *WsNet) Stop() {
	err := w.Engine.Stop(context.Background())
	if err != nil {
		logx.Error("ws net stop err", zap.Error(err))
		return
	}
}

func (w *WsNet) Run(ctx context.Context) {
	parse, err := url.Parse(w.wsAddr)
	if err != nil {
		logx.Panic("ws net parse add err", zap.String("addr", w.wsAddr), zap.Error(err))
		return
	}

	logx.Info("websocket server", zap.String("scheme", parse.Scheme),
		zap.String("host", parse.Host),
		zap.String("path", parse.Path),
		zap.String("port", parse.Port()),
		zap.String("query", parse.RawQuery),
		zap.String("fragment", parse.Fragment),
	)

	if parse.Scheme == "wss" {
		logx.Panic("暂不支持wss协议")
		return
	}

	err = gnet.Run(w, parse.Host,
		gnet.WithTicker(true),
		gnet.WithMulticore(true),
	)

	if err != nil {
		logx.Fatal("ws gnet run err", zap.Error(err))
	}
}

func New(addr string, handler flynet.EventHandler, opt *flynet.Options) *WsNet {
	res := &WsNet{
		wsAddr:  addr,
		opt:     opt,
		handler: handler,
		bufPool: pool.NewStdPool[*bytes.Buffer](func() *bytes.Buffer {
			return &bytes.Buffer{}
		}),
		framePool: opt.Pool,
	}
	return res
}
