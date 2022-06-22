package client

type Call struct {
	Seq           uint64      // 调用uuid
	ServiceMethod string      // 标准格式 "<service>.<method>"
	Args          interface{} // 调用方法入参
	Reply         interface{} // 调用方法返回
	Error         error       // 错误
	Done          chan *Call  // 方法调用完成的回调
}

// 调用结束时，会调用通知调用方。
func (call *Call) done() {
	call.Done <- call
}
