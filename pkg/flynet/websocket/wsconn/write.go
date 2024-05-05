package wsconn

import (
	"github.com/panjf2000/gnet/v2"
)

func (w *WsConn) Write(data []byte) (int, error) {
	frame := w.getFrame()
	defer w.freeFrame(frame)
	msg, err := frame.WriteBinary(data)
	if err != nil {
		return 0, err
	}
	return w.conn.Write(msg)
}

func (w *WsConn) AsyncWrite(data []byte) {

	frame := w.getFrame()
	msg, err := frame.WriteBinary(data)
	if err != nil {
		panic(err)
	}
	_ = w.conn.AsyncWrite(msg, func(c gnet.Conn, err error) error {
		w.freeFrame(frame)
		return err
	})
}
