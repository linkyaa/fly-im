package base

//
//import (
//	"github.com/linkyaa/fly-im/pkg/pool"
//	"runtime"
//	"time"
//)
//
//type (
//	Options struct {
//		NumDownLoop      int                 //处理下行消息的协程数,默认为CPU核数
//		AuthTimeout      time.Duration       //连接认证超时时间
//		KeepaliveTimeout time.Duration       //心跳超时时间
//		Pool             pool.Pooler[*Frame] //frame内存池
//	}
//
//	// AddrOptions 网关服务监听地址,如果对应字段不为空,则启动对应服务
//	AddrOptions struct {
//		WsAddr  string //Websocket服务监听地址
//		TcpAddr string //Tcp服务监听地址
//	}
//
//	Option func(opt *Options)
//)
//
//func WithNumDownLoop(num int) Option {
//	return func(opt *Options) {
//		opt.NumDownLoop = num
//	}
//}
//
//func WithKeepaliveTimeout(timeout time.Duration) Option {
//	return func(opt *Options) {
//		opt.KeepaliveTimeout = timeout
//	}
//}
//
//func WithAuthTimeout(timeout time.Duration) Option {
//	return func(opt *Options) {
//		opt.AuthTimeout = timeout
//	}
//}
//
//func WithPool(p pool.Pooler[*Frame]) Option {
//	return func(opt *Options) {
//		opt.Pool = p
//	}
//}
//
//func NewDefaultOptions() *Options {
//	return &Options{
//		NumDownLoop:      runtime.GOMAXPROCS(0),
//		AuthTimeout:      time.Second * 10,
//		KeepaliveTimeout: time.Minute * 5,
//		Pool: pool.NewStdPool[*Frame](func() any {
//			return &Frame{
//				FrameHeader: FrameHeader{},
//				FrameBody:   make([]byte, 0, 4096),
//			}
//		}),
//	}
//}
//
//func NewDefaultAddrOption() *AddrOptions {
//	return &AddrOptions{
//		WsAddr: "ws://127.0.0.1:5001/",
//	}
//}
