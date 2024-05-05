package flynet

import (
	"github.com/linkyaa/fly-im/pkg/flynet/frame"
	"github.com/linkyaa/fly-im/pkg/pool"
	"runtime"
	"time"
)

type (
	Options struct {
		NumDownLoop      int                         //处理下行消息的协程数,默认为CPU核数
		Pool             pool.Pooler[*frame.WsFrame] //frame池
		AuthTimeout      time.Duration               //建立好连接后,认证超时时间,默认20s
		KeepaliveTimeout time.Duration               //心跳超时时间,超时则关闭连接,当连接上有消息收发时,时间重置
	}

	AddrOption struct {
		WsAddr  string //Websocket服务监听地址
		TcpAddr string //Tcp服务监听地址
	}

	Option func(options *Options)
)

// WithNumDownLoop 处理下行消息的循环loop. TODO: 性能测试一下
func WithNumDownLoop(num int) Option {
	return func(options *Options) {
		options.NumDownLoop = num
	}
}

// WithKeepaliveTimeout connMgr管理的心跳超时配置,用于清除不活跃的连接
func WithKeepaliveTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.KeepaliveTimeout = timeout
	}
}

// WithAuthTimeout 鉴权的超时时间,.
// 这会给im的接入方带来问题, 或者鉴权就不做了, 由统一的鉴权网关去做.
func WithAuthTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.AuthTimeout = timeout
	}
}

func WithDownFramePool(framePool pool.Pooler[*frame.WsFrame]) Option {
	return func(options *Options) {
		options.Pool = framePool
	}
}

func NewDefaultOptions() *Options {
	return &Options{
		NumDownLoop:      runtime.GOMAXPROCS(0),
		AuthTimeout:      time.Second * 10,
		KeepaliveTimeout: time.Minute * 5,
		Pool: pool.NewStdPool[*frame.WsFrame](func() *frame.WsFrame {
			return frame.NewWsFrame(4096)
		}),
	}
}

func NewDefaultAddrOptions() *AddrOption {
	return &AddrOption{
		WsAddr:  "ws://0.0.0.0:8001/",
		TcpAddr: "",
	}
}
