package store

type (
	ConnContext interface {
		// DeviceId 设备ID
		DeviceId() int64

		// UserId 用户ID
		UserId()

		// DeviceType 设备类型
		DeviceType() uint8

		// ConnectTime 连接到服务器时间
		ConnectTime() int64

		// IsAuth 是否鉴权
		IsAuth() bool

		// Close 关闭连接
		Close() error
	}

	ConnMgr interface {
		// Add 添加一个新的conn,该conn未经过鉴权
		Add(ctx ConnContext)

		// GetByUid 通过uid获取
		GetByUid(uid int64, devices []ConnContext)

		GetByDeviceId(did int64) ConnContext

		// Remove 移除conn
		Remove(ctx ConnContext)
	}

	connMgr struct {
	}
)
