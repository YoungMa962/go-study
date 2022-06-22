package client

import (
	"codec"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

var ErrShutdown = errors.New("connection is shut down")
var _ io.Closer = (*Client)(nil) // 实现io.Closer 接口
type Client struct {
	cc       codec.Codec      // 消息的编解码器
	opt      *codec.Option    // 传输协议
	header   codec.Header     // 每个请求的消息头
	seq      uint64           // 用于给发送的请求编号，每个请求拥有唯一编号
	pending  map[uint64]*Call // 存储未处理完的请求，键是编号(seq),值是 Call 实例
	closing  bool             // client 用户主动关闭
	shutdown bool             // shutdown 置为 true 一般是有错误发生
	sending  sync.Mutex       // 互斥锁保证请求的有序发送，即防止出现多个请求报文混淆
	mu       sync.Mutex       // 互斥锁
}

type Result struct {
	client *Client
	err    error
}

//
// Dial 实现超时处理
// @Description: 创建连接
// @param network 网络协议
// @param address 地址
// @param opts 传输协议
// @return client
// @return err
//
func Dial(network, address string, opts ...*codec.Option) (client *Client, err error) {
	opt, err := parseOptions(opts...)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTimeout(network, address, opt.ConnectionTimeOut)
	if err != nil {
		return nil, err
	}
	// close the connection if client is nil
	defer func() {
		if client == nil {
			_ = conn.Close()
		}
	}()
	ch := make(chan Result)
	go NewClient(conn, opt, ch)

	// 如果为0，则为不设限制
	if opt.ConnectionTimeOut == 0 {
		result := <-ch
		return result.client, result.err
	}

	//select 语句阻塞等待最先返回数据的channel
	select {
	case <-time.After(opt.ConnectionTimeOut):
		return nil, fmt.Errorf("rpc client: connect timeout: expect within %s", opt.ConnectionTimeOut)
	case result := <-ch:
		return result.client, result.err
	}
}

//
// Go
// @Description: 异步调用，返回call实例ff
// @receiver client
// @param serviceMethod
// @param args
// @param reply
// @param done
// @return *Call
//
func (client *Client) Go(serviceMethod string, args, reply interface{}, done chan *Call) *Call {
	if done == nil {
		done = make(chan *Call, 10)
	} else if cap(done) == 0 {
		log.Panic("rpc client: done channel is unbuffered")
	}
	call := &Call{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Done:          done,
	}
	//time.Sleep(15 * time.Second)
	go client.send(call)
	return call
}

//
// Call
// @Description: 同步调用，等待call的结果并返回
// @receiver client
// @param serviceMethod
// @param args
// @param reply
// @return error
//
func (client *Client) Call(context context.Context, serviceMethod string, args, reply interface{}) error {
	call := client.Go(serviceMethod, args, reply, make(chan *Call, 1))
	select {
	case <-context.Done():
		client.removeCall(call.Seq)
		return errors.New("rpc client: call failed: call timeout")
	case call := <-call.Done:
		return call.Error
	}
}

func (client *Client) send(call *Call) {
	client.sending.Lock()
	defer client.sending.Unlock()

	seq, err := client.registerCall(call)
	if err != nil {
		call.Error = err
		call.done()
		return
	}

	// prepare request header
	client.header.ServiceMethod = call.ServiceMethod
	client.header.Seq = seq
	client.header.ResError = ""

	// encode and send the request
	if err := client.cc.Write(&client.header, call.Args); err != nil {
		call := client.removeCall(seq)
		// call may be nil, it usually means that Write partially failed,
		// client has received the response and handled
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

//
// parseOptions
// @Description: 解析传输协议
// @param opts
// @return *Option
// @return error
//
func parseOptions(opts ...*codec.Option) (*codec.Option, error) {
	// if opts is nil or pass nil as parameter
	if len(opts) == 0 || opts[0] == nil {
		return codec.DefaultOption, nil
	}
	if len(opts) != 1 {
		return nil, errors.New("number of options is more than 1")
	}
	opt := opts[0]
	opt.MagicNumber = codec.DefaultOption.MagicNumber
	if opt.CodecType == "" {
		opt.CodecType = codec.DefaultOption.CodecType
	}
	return opt, nil
}

func NewClient(conn net.Conn, opt *codec.Option, ch chan Result) {
	// 找到一个编解码器
	result := Result{client: nil, err: nil}
	codecFunc := codec.NewCodecFuncMap[opt.CodecType]
	if codecFunc == nil {
		err := fmt.Errorf("invalid codec type %s", opt.CodecType)
		log.Println("rpc client: codec error:", err)
		result.err = err
		ch <- result
		return
	}
	if err := json.NewEncoder(conn).Encode(opt); err != nil {
		log.Println("rpc client: options error: ", err)
		_ = conn.Close()
		result.err = err
		ch <- result
		return
	}
	result.client = newClientCodec(codecFunc(conn), opt)
	// time.Sleep(15 * time.Second)
	ch <- result
}

//
// newClientCodec
// @Description: 创建新的客户端实例，并开始接受消息
// @param cc
// @param opt
// @return *Client
//
func newClientCodec(cc codec.Codec, opt *codec.Option) *Client {
	client := &Client{
		seq:     1, // seq 从1开始，0代表无效
		cc:      cc,
		opt:     opt,
		pending: make(map[uint64]*Call),
	}
	go client.receive()
	return client
}

//
// receive
// @Description: 接收响应
// @receiver client
//
func (client *Client) receive() {
	var err error
	//不出错则一直处理
	for err == nil {
		var header codec.Header
		if err := client.cc.ReadHeader(&header); err != nil {
			break
		}
		// 通过 header.Seq 找到对应请求
		call := client.removeCall(header.Seq)
		switch {
		// 不存在 call
		case call == nil:
			err = client.cc.ReadBody(nil)
		// 返回请求结果有error
		case header.ResError != "":
			call.Error = fmt.Errorf(header.ResError)
			err = client.cc.ReadBody(nil)
			call.done()
		// 读取结果
		default:
			err = client.cc.ReadBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body " + err.Error())
			}
			call.done()
		}
	}
	client.terminateCalls(err)
}

func (client *Client) Close() error {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.closing {
		return ErrShutdown
	}
	client.closing = true
	return client.cc.Close()
}

// IsAvailable return true if the client does work
func (client *Client) IsAvailable() bool {
	client.mu.Lock()
	defer client.mu.Unlock()
	return !client.shutdown && !client.closing
}

//
// registerCall
// @Description: 将参数 call 添加到 client.pending 中，并更新 client.seq
// @receiver client
// @param call
// @return uint64
// @return error
//
func (client *Client) registerCall(call *Call) (uint64, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.closing || client.shutdown {
		return 0, ErrShutdown
	}
	call.Seq = client.seq
	client.pending[call.Seq] = call
	client.seq++
	return call.Seq, nil
}

//
// removeCall
// @Description: 根据 seq，从 client.pending 中移除对应的 call，并返回
// @receiver client
// @param seq
// @return *Call
//
func (client *Client) removeCall(seq uint64) *Call {
	client.mu.Lock()
	defer client.mu.Unlock()
	call := client.pending[seq]
	delete(client.pending, seq)
	return call
}

//
// terminateCalls
// @Description: 服务端或客户端发生错误时调用
// @receiver client
// @param err
//
func (client *Client) terminateCalls(err error) {
	client.sending.Lock()
	defer client.sending.Unlock()
	client.mu.Lock()
	defer client.mu.Unlock()
	//将 shutdown 设置为 true
	client.shutdown = true
	//将错误信息通知所有 pending 状态的 call
	for _, call := range client.pending {
		call.Error = err
		call.done()
	}
}
