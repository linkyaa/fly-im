package flynet

type (
	FrameHeader struct {
		//	TODO:没想好有什么类型的字段
	}
	Frame struct {
		FrameHeader
		FrameBody []byte
	}
)
