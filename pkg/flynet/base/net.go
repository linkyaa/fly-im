package base

import "net"

// flynet包作为长连接网关的网络引擎，是长连接网关的流量出入口
// 考虑长连接网关要承载海量连接，Go原生的网络库会使用大量的协程,所以采用基于事件驱动的网络库.
type (

	// Conn 事件循环中的conn
	Conn interface {
		//	============== 连接基本信息 ================

		// SetAuth 设置conn经过鉴权
		SetAuth(auth bool)
		// IsAuth 是否经过鉴权
		IsAuth() bool
		// RemoteAddr 对端地址
		RemoteAddr() net.Addr
		// LocalAddr 本地地址
		LocalAddr() net.Addr

		//DeviceId()int64

		// UserId 用户ID
		UserId() int64

		// ============== Read 接口 ==============

		// Read 从conn读取完整的数据包
		//Read(buf []byte) (int, error)

		// GetFrames 获取完整的frame,考虑到对接上层，还是返回frame比较好
		GetFrames() ([]*Frame, error)
		// ReleaseFrames 释放frame到内存池,释放之后frame将不可再用
		ReleaseFrames()

		// ============== Write 接口 ==============

		// Write 将frame写入conn,并发不安全,只能在事件循环中使用
		Write(frame *Frame) error
		// AsyncWrite 将frame写入conn,并发安全，异步的写入
		AsyncWrite(frame *Frame) error

		// ============== 控制接口 ==============

		// Close 关闭底层conn
		Close() error
	}
	// EventHandler 长连接网关的网络事件
	EventHandler interface {
		// OnData 收到完整数据包时触发
		OnData(conn Conn)

		// OnConnect 和服务器建立好连接时触发,对于ws协议来说，则是升级完成
		OnConnect(conn Conn)

		// OnClose 和服务器断开连接时触发
		OnClose(conn Conn, err error)
	}
)
