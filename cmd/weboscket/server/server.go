package main

import (
	"context"
	"github.com/alitto/pond"
	"github.com/gorilla/websocket"
	"github.com/linkyaa/fly-im/pkg/flynet"
	flyengine "github.com/linkyaa/fly-im/pkg/flynet/engine"
	"github.com/linkyaa/fly-im/pkg/logx"
	"github.com/linkyaa/fly-im/pkg/pool"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

type (
	echoDemo struct {
		bsPool        pool.Pooler[[]byte]
		connNum       int64
		requestQps    int64
		msgWorkerPool *pond.WorkerPool
	}
)

func (e *echoDemo) OnData(conn flynet.Conn) {
	var buf = e.bsPool.Get()

	for {
		res, size, done := conn.BytesReader(buf)

		if done {
			e.bsPool.Put(buf)
			return
		}

		buf = res
		e.msgWorkerPool.Submit(func() {
			conn.AsyncWrite(buf[:size])
			e.bsPool.Put(buf)
		})
		atomic.AddInt64(&e.requestQps, 1)
	}

}

func (e *echoDemo) OnConnect(conn flynet.Conn) {
	atomic.AddInt64(&e.connNum, 1)
}

func (e *echoDemo) OnClose(conn flynet.Conn, err error) {
	atomic.AddInt64(&e.connNum, -1)
}

func (e *echoDemo) monitor() {

	ticker := time.NewTicker(time.Second * 10)
	for {
		<-ticker.C
		logx.Info("monitor",
			zap.Int64("qps", atomic.SwapInt64(&e.requestQps, 0)/10),
			zap.Int64("conn", atomic.LoadInt64(&e.connNum)),
		)
	}

}

func runGnet(addr string) {
	logx.Info("start server ...")
	var ed = &echoDemo{
		bsPool: pool.NewStdPool[[]byte](func() []byte {
			return make([]byte, 4096)
		}),
		msgWorkerPool: pond.New(10000, 100000, pond.Strategy(pond.Lazy())), //懒惰回收
	}
	addrOption := flynet.NewDefaultAddrOptions()
	addrOption.WsAddr = "ws://0.0.0.0:9000/"
	engine := flyengine.NewEngine(ed, addrOption)
	go engine.Run(context.Background())
	go ed.monitor()
	waitDone()
	engine.Stop()
}

func main() {
	logx.SetLogLevel(zap.InfoLevel) //实现了动态变更的功能
	runGnet("9000")                 //行,先这样吧.
	//runGorilla(":9000")
}

func waitDone() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Kill, os.Interrupt)
	<-ch

	logx.Warn("server done .")
}

func runGorilla(addr string) {
	var cliNum int64 = 0
	var qps int64 = 0
	var upgrader = websocket.Upgrader{}
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			panic(err)
		}

		atomic.AddInt64(&cliNum, 1)
		go func() {
			defer atomic.AddInt64(&cliNum, -1)
			for {
				_, data, err := conn.ReadMessage()
				if err != nil {
					_ = conn.Close()
					return
				}

				err = conn.WriteMessage(websocket.BinaryMessage, data)
				if err != nil {
					_ = conn.Close()
					return
				}

				atomic.AddInt64(&qps, 1)
			}
		}()
	})

	go func() {
		ticker := time.NewTicker(time.Second * 10)
		for {
			<-ticker.C
			logx.Info("monitor",
				zap.Int64("qps", atomic.SwapInt64(&qps, 0)/10),
				zap.Int64("conn", atomic.LoadInt64(&cliNum)),
			)
		}
	}()
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
