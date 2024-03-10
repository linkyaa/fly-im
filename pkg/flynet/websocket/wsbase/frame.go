package wsbase

import (
	"encoding/binary"
	"github.com/gobwas/ws"
	"math/rand"
)

const (
	bit0 = 0x80
	bit1 = 0x40
	bit2 = 0x20
	bit3 = 0x10
	bit4 = 0x08
	bit5 = 0x04
	bit6 = 0x02
	bit7 = 0x01

	len7  = int64(125)
	len16 = int64(^(uint16(0)))
	len64 = int64(^(uint64(0)) >> 1)
)

// WsFrame 构造websocket协议的frame
type (
	WsFrame struct {
		ws.Header
		Payload                        []byte
		headerSize, payloadSize, total int
		header                         [ws.MaxHeaderSize]byte //最大头部
		rootBuf                        []byte
		//mask                           [4]byte
		//Frame                          ws.Frame
	}
)

func NewWsFrame(size int) *WsFrame {
	return &WsFrame{rootBuf: make([]byte, 0, size)}
}

func (w *WsFrame) WriteBinary(payload []byte) (data []byte, err error) {
	//1. mask payload
	w.Header.Fin = true
	w.Header.OpCode = ws.OpBinary
	w.Header.Length = int64(len(payload))
	w.payloadSize = len(payload)

	if true {
		binary.BigEndian.PutUint32(w.Mask[:], rand.Uint32())
		w.Masked = true
		ws.Cipher(payload, w.Mask, 0)
	} else {
		w.Masked = false
	}
	//w.Frame = ws.MaskFrameInPlaceWith(w.Frame, w.mask)
	//2. write header
	w.headerSize, err = WriteHeader(w.header[:], w.Header)
	if err != nil {
		return
	}
	//3. write payload
	w.tryGrow()
	copy(w.rootBuf[:w.headerSize], w.header[:w.headerSize])
	copy(w.rootBuf[w.headerSize:w.total], payload)
	data = w.rootBuf[:w.total]
	return
}

func (w *WsFrame) tryGrow() {
	total := w.headerSize + w.payloadSize
	if cap(w.rootBuf) < total {
		w.rootBuf = make([]byte, total)
	}

	w.total = total
}

// WriteHeader 要保证 bts 的len 足够
func WriteHeader(bts []byte, h ws.Header) (int, error) {
	// Make slice of bytes with capacity 14 that could hold any header.
	//bts := make([]byte, MaxHeaderSize)
	if len(bts) < ws.MaxHeaderSize {
		panic("bts len 小于 MaxHeaderSize")
	}

	if h.Fin {
		bts[0] |= bit0
	}
	bts[0] |= h.Rsv << 4
	bts[0] |= byte(h.OpCode)

	var n int
	switch {
	case h.Length <= len7:
		bts[1] = byte(h.Length)
		n = 2

	case h.Length <= len16:
		bts[1] = 126
		binary.BigEndian.PutUint16(bts[2:4], uint16(h.Length))
		n = 4

	case h.Length <= len64:
		bts[1] = 127
		binary.BigEndian.PutUint64(bts[2:10], uint64(h.Length))
		n = 10

	default:
		return 0, ws.ErrHeaderLengthUnexpected
	}

	if h.Masked {
		bts[1] |= bit0
		n += copy(bts[n:], h.Mask[:])
	}

	//_, err := w.Write(bts[:n])

	return n, nil
}
