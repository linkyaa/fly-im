package appruntime

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/linkyaa/fly-im/pkg/logx"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

type (
	Runner interface {
		// Run 启动服务
		Run(ctx context.Context)
		// Stop 停止服务
		Stop()
	}

	RunnerBuilder func() Runner

	option struct {
		appName        string
		beforeRun      func()
		afterRun       func()
		afterStop      func()
		waitStopSignal func()
	}

	Option func(opt *option)
)

// BeforeRun 在调用runner.Run之前执行
func BeforeRun(f func()) Option {
	return func(opt *option) {
		opt.beforeRun = f
	}
}

// AfterStop 在调用runner.Stop之后
func AfterStop(f func()) Option {
	return func(opt *option) {
		opt.afterStop = f
	}
}

// WaitStopSignal 等待退出信号
func WaitStopSignal(f func()) Option {
	return func(opt *option) {
		opt.waitStopSignal = f
	}
}

func newOptions(appName string) *option {
	return &option{
		appName: appName,
		beforeRun: func() {
			logx.Info("beforeRun:程序即将运行", zap.String("app", appName))
		},
		afterRun: func() {
			logx.Info("afterRun:程序运行中", zap.String("app", appName))
		},
		afterStop: func() {
			logx.Info("afterStop:程序正在退出", zap.String("app", appName))
		},
		waitStopSignal: func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Kill, os.Interrupt)
			sg := <-ch

			logx.Warn("waitStopSignal:收到退出信号",
				zap.String("signal", sg.String()),
				zap.String("app", appName),
			)
		},
	}
}

// NewApp 创建一个application.
// TODO: 解析配置文件
func NewApp(appName string, appOptions AppOptions, rb RunnerBuilder, opts ...Option) {

	//构建flag
	appOptions.AddFlags()
	errs := appOptions.Validate()
	if len(errs) != 0 {
		err := errors.Join(errs...)
		logx.Fatal(appName+" option validate err", zap.Error(err))
	}

	//添加基础flag
	initBasicFlag()

	flag.Parse()

	info, err := json.MarshalIndent(appOptions, "", " ")
	if err != nil {
		logx.Fatal("unmarshal option err", zap.Error(err))
	}
	//2. 输出flags
	fmt.Printf("options: %s\n", string(info))

	opt := newOptions(appName)

	for _, apply := range opts {
		apply(opt)
	}

	runner := rb()

	//3. beforeRun
	opt.beforeRun()

	//context是否可以由外部传入
	runner.Run(context.Background())

	opt.afterRun()

	opt.waitStopSignal()

	runner.Stop()
	opt.afterStop()
}
