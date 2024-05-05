package frame

import (
	"encoding/binary"
	"github.com/gobwas/ws"
)

var (
	EmptyHeader ws.Header
)

// ReadClientHeader reads a frame header from r.
// 先检查err, 在检查done
// size 表示已经读取的全部的size
func ReadClientHeader(h *ws.Header, bts []byte) (done bool, size int, err error) {
	// Make slice of bytes with capacity 12 that could hold any header.
	//
	// The maximum header size is 14, but due to the 2 hop reads,
	// after first hop that reads first 2 constant bytes, we could reuse 2 bytes.
	// So 14 - 2 = 12.
	//bts := make([]byte, 2, ws.MaxHeaderSize-2)

	// Prepare to hold first 2 bytes to choose size of next read.
	//_, err = io.ReadFull(r, bts)
	//if err != nil {
	//	return
	//}

	h.Fin = bts[0]&bit0 != 0
	h.Rsv = (bts[0] & 0x70) >> 4
	h.OpCode = ws.OpCode(bts[0] & 0x0f)

	var extra int //

	if bts[1]&bit0 != 0 {
		h.Masked = true
		extra += 4
	}

	length := bts[1] & 0x7f
	switch {
	case length < 126:
		h.Length = int64(length)

	case length == 126:
		extra += 2 //还有2字节作为实际的长度

	case length == 127:
		extra += 8 //还有8字节作为实际的长度

	default:
		err = ws.ErrHeaderLengthUnexpected
		return
	}

	//header解析完毕
	if extra == 0 {
		done = true
		return
	}

	size = 2 + extra // 2byte + extra

	//如果数据buf不够extra的需要,那么这个包没有解析完毕
	if len(bts) < size { //前2固定+extra长度
		done = false
		return
	}

	// Increase len of bts to extra bytes need to read.
	// Overwrite first 2 bytes that was read before.
	bts = bts[2 : 2+extra]

	switch {
	case length == 126:
		h.Length = int64(binary.BigEndian.Uint16(bts[:2]))
		bts = bts[2:]

	case length == 127:
		if bts[0]&0x80 != 0 {
			err = ws.ErrHeaderLengthMSB
			return
		}
		h.Length = int64(binary.BigEndian.Uint64(bts[:8]))
		bts = bts[8:]
	}

	if h.Masked {
		copy(h.Mask[:], bts)
	}

	done = true
	return
}
