package flynet

import (
	"fmt"
	"github.com/linkyaa/fly-im/pkg/flynet/base"
	tcpnet "github.com/linkyaa/fly-im/pkg/flynet/tcp"
	websocketnet "github.com/linkyaa/fly-im/pkg/flynet/websocket"
	"github.com/linkyaa/fly-im/pkg/pool"
)

/*
长连接网关服务的网络引擎,考虑到需要维持大量连接,因此采用基于事件驱动的框架,减少协程的开销.
TODO: 性能测试,协议完整性测试
*/

type (
	Engine struct {
		*base.Options
		ws      *websocketnet.WsNet
		tcp     *tcpnet.TcpNet
		handler base.EventHandler
		pool    pool.Pooler[*base.Frame]
	}
)

func (e *Engine) Run() {
	if e.ws != nil {
		//启动websocket server
		go e.ws.Run()
	}

	if e.tcp != nil {
		//启动tcp server
		go e.tcp.Run()
	}
}

func NewNetEngine(handler base.EventHandler, addr *base.AddrOptions, opts ...base.Option) *Engine {
	opt := base.NewDefaultOptions()
	for _, apply := range opts {
		apply(opt)
	}

	if addr.WsAddr == "" && addr.TcpAddr == "" {
		panic(fmt.Errorf("wsAddr && tcpAddr 不能同时为空,请指定监听地址"))
	}

	if opt.NumDownLoop == 0 {
		panic(fmt.Errorf("downLoop 的数量不能为0"))
	}

	res := &Engine{
		handler: handler,
		Options: opt,
		pool:    opt.Pool,
	}

	if addr.WsAddr != "" {
		res.ws = websocketnet.New(addr.WsAddr, handler, opt)
	}

	if addr.TcpAddr != "" {
		res.tcp = tcpnet.New()
	}

	return res

}
