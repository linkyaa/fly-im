package main

//import (
//	"errors"
//	"github.com/linkyaa/fly-im/pkg/logx"
//	"github.com/linkyaa/fly-im/pkg/wsclient"
//	"github.com/linkyaa/fly-im/pkg/wsclient/wsconn"
//	"go.uber.org/zap"
//	"sync/atomic"
//	"time"
//)
//
//type (
//	clientEngine struct {
//	}
//)
//
//var (
//	failConn int64
//)
//
//func (c *clientEngine) OnOpen(conn *wsconn.WsConn) {
//	logx.Debug("ws 协议升级完成")
//}
//
//func (c *clientEngine) OnClose(conn *wsconn.WsConn, err error) {
//	if err == nil {
//		err = errors.New("unknown")
//	}
//	logx.Error("on close", zap.String("reason", err.Error()))
//	atomic.AddInt64(&failConn, 1)
//}
//
//func (c *clientEngine) OnData(conn *wsconn.WsConn) {
//	frames := conn.GetFrames()
//	if len(frames) == 0 {
//		return
//	}
//
//	logx.Info("frame", zap.Int("len", len(frames)),
//		zap.String("msg", string(frames[0].Payload)))
//	conn.ReleaseFrames()
//}
//
//
//
//func runWs(addr string) {
//	var engine = &clientEngine{}
//
//	client := wsclient.New(engine)
//	client.Run()
//
//	logx.Info("new ws client")
//	var conn *wsconn.WsConn
//	var start = time.Now().UTC().UnixMilli()
//
//	conn = client.DialWs(addr, nil)
//	logx.Warn("connect to server done",
//		zap.Int64("fail", atomic.LoadInt64(&failConn)),
//		zap.Int64("time(ms)", time.Now().UTC().UnixMilli()-start),
//	)
//
//	time.Sleep(time.Second)
//
//	handleErr(conn.Write([]byte("hello world")))
//
//	logx.Debug("发送signIn消息成功")
//	select {}
//}
//
//func handleErr(err error) {
//	if err != nil {
//		panic(err)
//	}
//}
//
//func main() {
//	go runWs("ws://localhost:8001/ws?name=ly#100")
//	select {}
//}
