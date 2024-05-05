package flynet

import (
	"net"
)

// flynet 作为长连接网关网络引擎,是长连接网关的流量出入口.
type (
	// Conn flynet event上层封装的conn
	// TODO: 需要整理一下哪些是有并发安全要求的
	Conn interface {
		// ================================== 连接的基本信息 ==================================

		// ConnectTime conn接入服务器时间
		ConnectTime() int64

		// SetAuth 设置conn是否经过认证,并发安全
		SetAuth()
		// IsAuth 是否经过认证
		IsAuth() bool
		// RemoteAddr returns the remote network address.
		RemoteAddr() net.Addr
		// LocalAddr returns the local network address.
		LocalAddr() net.Addr

		// DeviceId 设备ID
		DeviceId() int64

		// SetDeviceId 设置设备ID,在用户登录成功后
		SetDeviceId(deviceId int64)

		// UserId 用户ID
		UserId() int64

		// SetUserId 设置用户ID,在用户登录成功后
		SetUserId(uid int64)

		// DeviceType 设备类型
		DeviceType() uint8

		// SetDeviceType 设置设备类型,在用户登录成功后
		SetDeviceType(deviceType uint8)

		// Id 这个ID只用于当前节点
		//Id() int64

		// ================================== Read 接口 ==================================

		// BytesReader 从conn读取完整的数据包到buf
		// TODO: 提供更友好的接口,把数据拆为多次提供就行. 然后提供这个接口 Read([]byte)(int,error)
		// buf: 从conn读取数据到buf,数据通过res返回,如果buf容量不够,内部新建一个buf.
		// res: 从conn读取的bytes
		// done:是否已经从conn读取完毕
		// err: 从conn读取发生error
		// size:本次从读取的byte的
		BytesReader(buf []byte) (res []byte, size int, done bool)

		// ================================== Write 接口 ==================================

		// Write 将data打包为数据从conn发送出去
		Write(data []byte) (int, error) //并发不安全的,需要加锁
		AsyncWrite(data []byte)         //异步的并发安全的写入数据

		// ================================== 控制接口 ==================================

		// Close 关闭连接
		Close() error
	}

	// EventHandler 长连接网关网络事件
	EventHandler interface {
		// OnData 收到完整的数据包
		OnData(conn Conn)
		// OnConnect 和服务器建立连接,处理auth超时
		OnConnect(conn Conn)
		// OnClose 和服务器断开连接
		OnClose(conn Conn, err error) // 需要封装Err
	}
)
