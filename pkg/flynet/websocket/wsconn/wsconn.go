package wsconn

import (
	"bytes"
	"github.com/linkyaa/fly-im/pkg/flynet/frame"
	"github.com/linkyaa/fly-im/pkg/pool"
	"github.com/panjf2000/gnet/v2"
	"io"
	"net"
	"sync/atomic"
	"time"
)

const (
	upgrading = iota
	completed
)

type (
	WsConn struct {
		status     status                      //conn状态
		conn       gnet.Conn                   //底层连接
		buf        *bytes.Buffer               //用于升级的buf,用完回收
		bufPool    pool.Pooler[*bytes.Buffer]  //buf内存池
		cacheFrame *frame.WsFrame              //未完整的frame
		framePool  pool.Pooler[*frame.WsFrame] //frame pool

		// ==== info ====
		auth        int64    //是否经过鉴权
		closed      int64    //是否关闭
		connectTime int64    //连接到服务器的时间,单位s
		deviceId    int64    //设备ID,登录后设置
		userId      int64    //用户ID,登录后设置
		deviceType  int32    //设备类型,登录后设置
		remoteAddr  net.Addr //conn的地址

		// ==== read ====
		availableFramesIndex int              //frames已读的index
		availableFrames      []*frame.WsFrame //完整的frame

		// ==== write ====
	}

	//websocket升级的reader/writer
	wsRW struct {
		io.Reader
		io.Writer
	}

	status uint8
)

func (w *WsConn) BytesReader(buf []byte) ([]byte, int, bool) {
	curIndex := w.availableFramesIndex
	if curIndex > len(w.availableFrames)-1 {
		w.freeFrames()
		return nil, 0, true
	}

	buf = append(buf[:0], w.availableFrames[curIndex].Payload...)
	w.availableFramesIndex++
	return buf, len(buf), false
}

// Close 可以重复关闭吗,TODO：应该统一gnet进行关闭,关闭的触发应该只有gnet的OnClose
func (w *WsConn) Close() error {
	//close 后, 和gnet相关的都不能用
	//因此对于异步任务,主要是写操作都要判断是否关闭,还要注意info相关的东西
	if atomic.CompareAndSwapInt64(&w.closed, 0, 1) {
		return w.conn.Close()
	}
	return nil
}

func (w *WsConn) growBuf(size int) {
	if w.buf == nil {
		w.buf = w.bufPool.Get()
	}
	w.buf.Grow(size)
}

func (w *WsConn) freeBuf() {
	w.buf.Reset()
	w.bufPool.Put(w.buf)
	w.buf = nil
}

func (w *WsConn) getFrame() *frame.WsFrame {
	return w.framePool.Get()
}

func (w *WsConn) freeFrame(f *frame.WsFrame) {
	f.Header = frame.EmptyHeader
	w.framePool.Put(f)
}

func (w *WsConn) freeFrames() {
	for _, f := range w.availableFrames {
		w.freeFrame(f)
	}
	w.availableFramesIndex = 0
	w.availableFrames = w.availableFrames[:0]
}

func NewWsConn(c gnet.Conn, bufPool pool.Pooler[*bytes.Buffer], framePool pool.Pooler[*frame.WsFrame]) *WsConn {
	res := &WsConn{
		status:          upgrading,
		conn:            c,
		bufPool:         bufPool,
		remoteAddr:      c.RemoteAddr(),
		framePool:       framePool,
		availableFrames: make([]*frame.WsFrame, 0),
		connectTime:     time.Now().UTC().Unix(),
	}
	return res
}
