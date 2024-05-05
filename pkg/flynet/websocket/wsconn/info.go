package wsconn

import (
	"net"
	"sync/atomic"
)

func (w *WsConn) IsAuth() bool {
	return atomic.LoadInt64(&w.auth) == 1
}

func (w *WsConn) SetDeviceId(deviceId int64) {
	atomic.StoreInt64(&w.deviceId, deviceId)
}

func (w *WsConn) SetUserId(uid int64) {
	atomic.StoreInt64(&w.userId, uid)
}

func (w *WsConn) SetDeviceType(deviceType uint8) {
	atomic.StoreInt32(&w.deviceType, int32(deviceType))
}

func (w *WsConn) SetAuth() {
	atomic.StoreInt64(&w.auth, 1)
}

func (w *WsConn) DeviceId() int64 {
	return atomic.LoadInt64(&w.deviceId)
}

func (w *WsConn) UserId() int64 {
	return atomic.LoadInt64(&w.userId)
}

func (w *WsConn) DeviceType() uint8 {
	return uint8(atomic.LoadInt32(&w.deviceType))
}

func (w *WsConn) ConnectTime() int64 {
	return w.connectTime
}

func (w *WsConn) RemoteAddr() net.Addr {
	return w.conn.RemoteAddr()
}

func (w *WsConn) LocalAddr() net.Addr {
	//TODO implement me
	panic("implement me")
}
