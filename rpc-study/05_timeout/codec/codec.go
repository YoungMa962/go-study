package codec

import "io"

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json" // not implemented
)

type NewCodecFunc func(io.ReadWriteCloser) Codec
type Type string

var NewCodecFuncMap map[Type]NewCodecFunc

type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
