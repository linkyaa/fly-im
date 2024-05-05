package flynet

import (
	"errors"
	"net/http"
)

//定义错误类型:wrap/unwrap/withStack

var (
	ErrWsHeaderToLarge      = errors.New("ws协议header超过消息")
	ErrConnRead             = errors.New("从conn读取buffer失败")
	ErrReadDataInconsistent = errors.New("读取数据不一致")
	ErrWsPeerClose          = errors.New("ws协议对端关闭,opClose")
)

const (
	WsMaxHeaderSize = http.DefaultMaxHeaderBytes
)
