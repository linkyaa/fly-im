package wsconn

import (
	"bytes"
	"github.com/gobwas/ws"
	"github.com/linkyaa/fly-im/pkg/flynet"
	"github.com/linkyaa/fly-im/pkg/flynet/frame"
	"github.com/linkyaa/fly-im/pkg/logx"
	"github.com/pkg/errors"
	"io"
)

// TryUpgrade 尝试升级
// bool: 升级是否完成
// err: 升级发生错误
func (w *WsConn) TryUpgrade() (bool, error) {
	if w.status == upgrading {
		return w.preUpgrade()
	}
	return true, nil
}

// ws升级预处理
func (w *WsConn) preUpgrade() (ok bool, err error) {
	size := w.conn.InboundBuffered()
	w.growBuf(size)

	if w.buf.Len() > flynet.WsMaxHeaderSize {
		return false, errors.WithStack(flynet.ErrWsHeaderToLarge)
	}

	readSize, err := w.buf.ReadFrom(w.conn)
	if err != nil {
		return false, errors.WithStack(flynet.ErrConnRead)
	}

	if readSize != int64(size) {
		return false, errors.WithStack(flynet.ErrReadDataInconsistent)
	}

	return w.upgrade()
}

// ws协议升级
func (w *WsConn) upgrade() (ok bool, err error) {
	if logx.EnableDebug() {
		logx.Debug("ws upgrade conn")
	}
	//1. 处理升级
	upgrader := &ws.Upgrader{}
	reader := bytes.NewReader(w.buf.Bytes())
	rw := &wsRW{
		Reader: reader,
		Writer: w.conn,
	}

	_, err = upgrader.Upgrade(rw)

	//检查升级是否完成了. 在upgrade过程中, 如果reader数据不够
	//则返回eof
	// || err == io.ErrUnexpectedEOF
	if err == io.EOF {
		return false, nil
	}

	if err != nil {
		return false, errors.WithStack(err)
	}

	//buffer中还存在为消费的数据,不符合http协议的预期
	if reader.Len() > 0 {
		return false, errors.WithStack(flynet.ErrReadDataInconsistent)
	}

	w.freeBuf()
	w.status = completed
	if logx.EnableDebug() {
		logx.Debug("ws upgrade done")
	}
	return true, nil
}

func (w *WsConn) DecodeFrame() (bool, error) {
	err := w.decodeFrame()
	//TODO：需要过滤 data 和 其他的类型,暂时放在read的时候处理。
	return len(w.availableFrames) != 0, err
}

func (w *WsConn) decodeFrame() error {
	size := w.conn.InboundBuffered()
	w.growBuf(size)
	readSize, err := w.buf.ReadFrom(w.conn)
	if err != nil {
		return errors.WithStack(err)
	}

	if readSize != int64(size) {
		return errors.WithStack(flynet.ErrReadDataInconsistent)
	}

	for {
		if w.cacheFrame != nil {
			//如果不是初始帧
			length := w.cacheFrame.Header.Length
			if length > int64(w.buf.Len()) {
				return nil
			}

			w.cacheFrame.Payload = append(w.cacheFrame.Payload, w.buf.Next(int(length))...)
			if w.cacheFrame.Header.Masked {
				ws.Cipher(w.cacheFrame.Payload, w.cacheFrame.Header.Mask, 0)
			}

			//TODO: 只处binary理类型的消息
			switch w.cacheFrame.OpCode {
			case ws.OpBinary:
				w.availableFrames = append(w.availableFrames, w.cacheFrame)
			case ws.OpClose:
				//TODO: 对于opClose, 应该降低一下日志级别, 不需要stack也行
				w.freeFrame(w.cacheFrame)
				return errors.WithStack(flynet.ErrWsPeerClose)
			default:
				w.freeFrame(w.cacheFrame)
			}
			w.cacheFrame = nil
		}

		if w.buf.Len() == 0 {
			w.freeBuf()
			return nil
		}

		//数据不够解析 header
		if w.buf.Len() < ws.MinHeaderSize {
			return nil
		}

		//如果buf残留未读取的new frame,那么先解析header
		f := w.getFrame()
		header := &f.Header
		done, hds, err := frame.ReadClientHeader(header, w.buf.Bytes())
		if err != nil {
			w.freeFrame(f)
			return err
		}

		//头部数据不完整
		if !done {
			return nil
		}

		w.cacheFrame = f
		w.cacheFrame.Payload = w.cacheFrame.Payload[:0]
		w.buf.Next(hds)
	}
}
