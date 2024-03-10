package main

import (
	"github.com/gobwas/ws"
	"github.com/linkyaa/fly-im/pkg/flynet"
	"github.com/linkyaa/fly-im/pkg/flynet/base"
	. "github.com/linkyaa/fly-im/pkg/logx"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"time"
)

type (
	eventHandle struct {
	}
)

func (e *eventHandle) OnData(conn base.Conn) {
	//TODO:多协议的支持,这里还不清楚怎么处理,也许判断协议类型.
	frames := conn.GetWsFrames()
	for _, frame := range frames {
		if frame.OpCode == ws.OpClose {
			Logger.Info("op close")
			_ = conn.Close()
			break
		}
		Logger.Info(string(frame.Payload), zap.Int("code", int(frame.OpCode)))
		err := conn.Write(&base.Frame{
			FrameHeader: base.FrameHeader{},
			FrameBody:   frame.Payload,
		})
		if err != nil {
			Logger.Error("write err", zap.Error(err))
			_ = conn.Close()
			break
		}
	}
	conn.ReleaseFrames()
}

func (e *eventHandle) OnConnect(conn base.Conn) {
	Logger.Info("on connect")
}

func (e *eventHandle) OnClose(conn base.Conn, err error) {
	Logger.Info("on close", zap.Error(err))
}

func main() {
	var event = &eventHandle{}
	engine := flynet.NewNetEngine(event, base.NewDefaultAddrOption())
	engine.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	engine.Stop()
	//TODO:协调stop.
	time.Sleep(time.Second)
	Logger.Info("stop server ...")
}
