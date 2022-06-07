package geerpc

import "codec"

const MagicNumber = 0x3bef5c

// Option
// @Description: | Option{MagicNumber: xxx, CodecType: xxx} | Header{ServiceMethod ...} | Body interface{} |
// 				 | <------      固定 JSON 编码      ------>  | <-------   编码方式由 CodeType 决定   ------->  |
//
type Option struct {
	MagicNumber int        // MagicNumber marks this's a geerpc request
	CodecType   codec.Type // client may choose different Codec to encode body//
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:   codec.GobType,
}
