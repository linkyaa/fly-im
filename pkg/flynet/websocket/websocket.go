package websocket_net

import (
	"github.com/panjf2000/gnet/v2"
	"time"
)

type (
	WsNet struct {
		gnet.Engine
		wsAddr string
	}
)

func (w *WsNet) OnBoot(eng gnet.Engine) (action gnet.Action) {
	w.Engine = eng
	return gnet.None
}

func (w *WsNet) OnShutdown(eng gnet.Engine) {
	//TODO implement me
	panic("implement me")
}

func (w *WsNet) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	//TODO implement me
	panic("implement me")
}

func (w *WsNet) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	//TODO implement me
	panic("implement me")
}

func (w *WsNet) OnTraffic(c gnet.Conn) (action gnet.Action) {
	//TODO implement me
	panic("implement me")
}

func (w *WsNet) OnTick() (delay time.Duration, action gnet.Action) {
	//TODO implement me
	panic("implement me")
}

func New() *WsNet {
	res := &WsNet{}
	return res
}
