package option

import (
	"flag"
	"fmt"
	"github.com/linkyaa/fly-im/pkg/appruntime"
)

var (
	_ appruntime.AppOptions = (*GRPCServerOption)(nil)
	_ appruntime.AppOptions = (*GRPCClientOption)(nil)
)

type (
	GRPCServerOption struct {
		BindAddress string `json:"bind_address"`
		BindPort    int    `json:"bind_port"`
		MaxMsgSize  int    `json:"max_msg_size"`
	}

	GRPCClientOption struct {
		ClientCert    string `json:"client_cert"`
		ServerAddress string `json:"server_address"`
	}
)

func (g *GRPCClientOption) AddFlags() {
	flag.StringVar(&g.ServerAddress, "grpc.server_address", g.ServerAddress, "grpc客户端连接到服务器的地址")
	flag.StringVar(&g.ClientCert, "grpc.client_cert", g.ClientCert, "grpc客户端证书")
}

func (g *GRPCClientOption) Validate() []error {
	if g.ServerAddress == "" {
		return []error{fmt.Errorf("GRPCClientOption.ServerAddress 服务器的地址不能为空")}
	}
	return nil
}

func (g *GRPCServerOption) Validate() []error {
	var errors []error
	if g.BindPort < 0 || g.BindPort > 65535 {
		errors = append(errors, fmt.Errorf("grpc 绑定非法端口 %d, 端口范围在 0 - 65535", g.BindPort))
	}
	return errors
}

func (g *GRPCServerOption) AddFlags() {
	flag.StringVar(&g.BindAddress, "grpc.bind_address", g.BindAddress,
		"grpc 服务监听地址")

	flag.IntVar(&g.BindPort, "grpc.bind_port", g.BindPort, "grpc 服务监听端口")

	flag.IntVar(&g.MaxMsgSize, "grpc.max_msg_size", g.MaxMsgSize, "grpc 服务接受消息的最大大小")
}

func NewGRPCOption() *GRPCServerOption {
	return &GRPCServerOption{
		BindAddress: "0.0.0.0",
		BindPort:    8100,
		MaxMsgSize:  4 << 20,
	}
}

func NewGRPCClientOption() *GRPCClientOption {
	return &GRPCClientOption{
		ClientCert:    "",
		ServerAddress: "127.0.0.1:8001",
	}
}
