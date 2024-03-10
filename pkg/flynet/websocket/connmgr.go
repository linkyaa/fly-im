package websocket_net

import (
	"github.com/linkyaa/fly-im/pkg/ring"
	"time"
)

// TODO:可以考虑增加连接上限
type (
	//连接管理:
	//1. 连接超时
	connMgr struct {
		upgradeTimeout int64 //从tcp升级到ws的超时时间,单位5s
		queue          *ring.Queue[*wsConn]
		dataIn         chan *wsConn
	}
)

func newConnMgr() *connMgr {
	res := &connMgr{
		upgradeTimeout: 5, //
		queue:          ring.New[*wsConn](),
		dataIn:         make(chan *wsConn, 1e5),
	}

	go res.startTicker()

	return res
}

func (cm *connMgr) startTicker() {
	ticker := time.NewTicker(time.Second * time.Duration(cm.upgradeTimeout))
	for {
		select {
		case <-ticker.C:
			if cm.queue.Length() == 0 {
				continue
			}
			//检查连接
			cm.check()
		case newConn := <-cm.dataIn:
			cm.queue.Add(newConn)
		}

	}
}

func (cm *connMgr) check() {
	now := time.Now().UTC().Unix()
	for cm.queue.Length() != 0 {
		conn := cm.queue.Peek()
		if conn.connectTime+cm.upgradeTimeout < now {
			//没有超时
			return
		}
		//TODO:增加认证超时
		if conn.status == upgrading {
			_ = conn.conn.Close() //不直接调用Close,并发不安全
		}
		cm.queue.Remove()
	}
}
func (cm *connMgr) add(conn *wsConn) {
	cm.dataIn <- conn
}
