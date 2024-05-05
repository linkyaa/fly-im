package gateway

import (
	"context"
	gatewayoption "github.com/linkyaa/fly-im/internal/gateway/option"
	"github.com/linkyaa/fly-im/pkg/appruntime"
	"github.com/linkyaa/fly-im/pkg/logx"
)

/*
gateway服务的主要任务：
1. 协议支持
2. 流量出入口
3. 消息转发
4. 长连接管理
*/

type (
	// Gateway 长连接网关服务
	server struct {
	}
)

func (s *server) Run(ctx context.Context) {
	logx.Info("impl me")
}

func (s *server) Stop() {
	logx.Info("impl me")
}

func NewServer(opt *gatewayoption.Options) appruntime.Runner {
	res := &server{}

	return res
}
