package gatewayoption

import (
	"flag"
	"github.com/linkyaa/fly-im/pkg/appruntime"
	"github.com/linkyaa/fly-im/pkg/option"
)

var (
	_ appruntime.AppOptions = (*Options)(nil)
)

type (
	Options struct {
		ConnAuthTimeout int64                    `json:"gw.conn_auth_timeout"`
		WsServerAddr    string                   `json:"gw.ws_server_addr"`
		GatewayCode     string                   `json:"gw.gateway_code"`
		Server          *option.GRPCServerOption `json:"server"`
		Log             *option.LogOption        `json:"log"`
		//Version         string                   `json:"gw.version"`
	}
)

func (o *Options) AddFlags() {
	o.Server.AddFlags()
	o.Log.AddFlags()
	flag.Int64Var(&o.ConnAuthTimeout, "gw.conn_auth_timeout", o.ConnAuthTimeout, "连接认证超时时间")
	flag.StringVar(&o.WsServerAddr, "gw.ws_server_addr", o.WsServerAddr, "ws服务监听地址")
	flag.StringVar(&o.GatewayCode, "gw.gateway_code", o.GatewayCode, "网关节点唯一码")
	//flag.StringVar(&o.Version, "gw.version", o.Version, "sdk版本")
}

func (o *Options) Validate() []error {
	var errs []error
	errs = append(errs, o.Log.Validate()...)
	errs = append(errs, o.Server.Validate()...)
	return errs
}

func NewOptions() *Options {
	//TODO: 设置gw code
	return &Options{
		WsServerAddr:    "ws://0.0.0.0:8001",
		ConnAuthTimeout: 10,
		GatewayCode:     "",
		Server:          option.NewGRPCOption(),
		Log:             option.NewLogOption(),
	}
}
