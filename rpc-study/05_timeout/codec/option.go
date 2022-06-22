package codec

import "time"

const MagicNumber = 0x3bef5c

// Option
// @Description: | Option{MagicNumber: xxx, CodecType: xxx} | Header{ServiceMethod ...} | Body interface{} |
// 				 | <------      固定 编码      ------>  | <-------   编码方式由 CodeType 决定   ------->  |
//
type Option struct {
	MagicNumber       int           // 用于校验
	CodecType         Type          // 编码方式
	ConnectionTimeOut time.Duration // 连接超时时间
	HandlerTimeOut    time.Duration // 处理超时时间
}

var DefaultOption = &Option{
	MagicNumber:       MagicNumber,
	CodecType:         GobType,
	ConnectionTimeOut: 10 * time.Second,
	HandlerTimeOut:    0 * time.Second,
}
