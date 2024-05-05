package tcp_net

import (
	"context"
	"github.com/linkyaa/fly-im/pkg/flynet"
)

type (
	TcpNet struct {
	}
)

func (t *TcpNet) Run(ctx context.Context) {

}

func (t *TcpNet) Stop() {

}

func New(addr string, handler flynet.EventHandler, opt *flynet.Options) *TcpNet {
	res := &TcpNet{}
	return res
}
