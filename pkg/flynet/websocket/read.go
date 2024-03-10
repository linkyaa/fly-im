package websocket_net

import (
	"bytes"
	"errors"
	"github.com/gobwas/ws"
	"github.com/linkyaa/fly-im/pkg/flynet/websocket/wsbase"
	"io"
	"net/http"
)

func (w *wsConn) procUpgrade() (ok bool, err error) {
	size := w.conn.InboundBuffered()
	w.bufGrow(size)
	if w.buf.Len() > http.DefaultMaxHeaderBytes {
		return false, errors.New("升级请求buf过大")
	}
	readSize, err := w.buf.ReadFrom(w.conn)
	if err != nil {
		return false, err
	}

	if readSize != int64(size) {
		return false, errors.New("conn 读取数据错误")
	}

	return w.upgrade()
}

func (w *wsConn) upgrade() (bool, error) {
	//1. 处理升级
	upgrade := &ws.Upgrader{}

	readerBuf := w.buf.Bytes()
	reader := bytes.NewReader(readerBuf)
	var rw = &upgrader{
		Reader: reader,
		Writer: w.conn,
	}

	_, err := upgrade.Upgrade(rw)
	//升级没有完成
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	remain := reader.Len()
	if remain > 0 {
		//TODO:按照Http的特点，在没有响应前，客户端应该不会发消息,所以检测到skip大于0,应该是不正常的状态.先关闭了.
		return false, errors.New("数据有残留,升级失败")
	}

	w.freeBuf()
	w.status = upgraded
	return true, nil
}

func (w *wsConn) procWsFrames() error {
	size := w.conn.InboundBuffered()
	w.bufGrow(size)
	readSize, err := w.buf.ReadFrom(w.conn)
	if err != nil {
		return err
	}

	if readSize != int64(size) {
		return errors.New("frame read size err")
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

			w.completeFrames = append(w.completeFrames, w.cacheFrame)
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
		frame := w.getFrame()
		header := &frame.Header
		done, hds, err := wsbase.ReadHeader(header, w.buf.Bytes())
		if err != nil {
			w.freeFrame(frame)
			return err
		}

		//头部数据不完整
		if !done {
			return nil
		}

		w.cacheFrame = frame
		w.cacheFrame.Payload = w.cacheFrame.Payload[:0]
		w.buf.Next(hds)
	}
}

func (w *wsConn) getCompleteFrames() []*wsbase.WsFrame {
	return w.completeFrames
}

func (w *wsConn) freeFrames() {
	for _, frame := range w.completeFrames {
		w.freeFrame(frame)
	}

	w.completeFrames = w.completeFrames[:0]

}
