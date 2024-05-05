package main

import (
	"github.com/linkyaa/fly-im/internal/gateway"
	gatewayoption "github.com/linkyaa/fly-im/internal/gateway/option"
	"github.com/linkyaa/fly-im/pkg/appruntime"
	"github.com/linkyaa/fly-im/pkg/logx"
)

func main() {
	options := gatewayoption.NewOptions()

	appruntime.NewApp("gateway", options, func() appruntime.Runner {
		runner := gateway.NewServer(options)

		logx.SetLogLevel(options.Log.ZapLevel)
		return runner
	})
}
