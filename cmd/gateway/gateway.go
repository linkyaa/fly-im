package main

import (
	"fmt"
	"github.com/linkyaa/fly-im/pkg/flynet"
	"github.com/linkyaa/fly-im/pkg/flynet/base"
	. "github.com/linkyaa/fly-im/pkg/logx"
)

func main() {
	Logger.Debug("hello world!!!")
	fmt.Print()
	flynet.NewNetEngine(nil, base.NewDefaultAddrOption())
}
