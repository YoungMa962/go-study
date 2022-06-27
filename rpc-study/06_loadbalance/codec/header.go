package codec

type Header struct {
	ServiceMethod string // 请求方法名
	Seq           uint64 // 请求uuid
	ResError      string // 服务端返回错误，客户端初始化为空
}

func NewHeader(serviceMethod string, seq uint64, resError string) *Header {
	return &Header{ServiceMethod: serviceMethod, Seq: seq, ResError: resError}
}
