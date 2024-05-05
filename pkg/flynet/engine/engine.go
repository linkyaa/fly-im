package flyengine

import (
	"context"
	"github.com/linkyaa/fly-im/pkg/flynet"
	tcp_net "github.com/linkyaa/fly-im/pkg/flynet/tcp"
	websocket_net "github.com/linkyaa/fly-im/pkg/flynet/websocket"
	"github.com/linkyaa/fly-im/pkg/logx"
)

type (
	//Engine 关注于net的控制.
	Engine struct {
		handler flynet.EventHandler //基于事件循环的event handle
		opt     *flynet.Options
		addr    *flynet.AddrOption

		//net
		wsNet  *websocket_net.WsNet
		tcpNet *tcp_net.TcpNet
	}
)

// Run 启动engine,可以考虑给run返回err处理
func (e *Engine) Run(ctx context.Context) {
	if e.addr.WsAddr == "" && e.addr.TcpAddr == "" {
		//Panic吧，或者unHealth也行.
		logx.Panic("engine addr 参数错误,地址不能全部为空")
	}

	if e.addr.WsAddr != "" {
		//启动ws net
		e.wsNet = websocket_net.New(e.addr.WsAddr, e.handler, e.opt)
		go e.wsNet.Run(ctx)
	}

	if e.addr.TcpAddr != "" {
		//启动tcp net
		e.tcpNet = tcp_net.New(e.addr.TcpAddr, e.handler, e.opt)
		go e.tcpNet.Run(ctx)
	}

	//等待退出信号
	<-ctx.Done()
}

// Stop 停止engine
func (e *Engine) Stop() {
	if e.wsNet != nil {
		e.wsNet.Stop()
	}
	if e.tcpNet != nil {
		e.tcpNet.Stop()
	}
}

// NewEngine 创建net engine
// [handler]
// [addr]
func NewEngine(handler flynet.EventHandler, addr *flynet.AddrOption, opts ...flynet.Option) *Engine {
	opt := flynet.NewDefaultOptions()

	for _, f := range opts {
		f(opt)
	}

	//TODO: 检查参数合法性

	res := &Engine{
		handler: handler,
		opt:     opt,
		addr:    addr,
	}

	return res
}
