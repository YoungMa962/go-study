package geerpc

import (
	"codec"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

var defaultServer = NewServer()

func Accept(lis net.Listener) { defaultServer.Accept(lis) }

func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
			return
		}
		go server.ServeConn(conn)
	}
}

func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	var opt Option
	// parse income message to option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: options error: ", err)
		return
	}
	// check MagicNumber
	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server: invalid magic number %x", opt.MagicNumber)
		return
	}

	codec := codec.NewCodecFuncMap[opt.CodecType]
	server.serveCodec(codec(conn))
}

// invalidRequest is a placeholder for response argValue when error occurs
var invalidRequest = struct{}{}

// readRequest		读取请求
// handleRequest	处理请求
// sendResponse 	回复请求
func (server *Server) serveCodec(cc codec.Codec) {
	sending := new(sync.Mutex) // make sure to send a complete response
	wg := new(sync.WaitGroup)  // wait until all request are handled
	for {
		req, err := server.readRequest(cc)
		if err != nil {
			if req == nil {
				break // it's not possible to recover, so close the connection
			}
			req.header.ResError = err.Error()
			server.sendResponse(cc, req.header, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go server.handleRequest(cc, req, sending, wg)
	}
	wg.Wait()
	_ = cc.Close()
}

// request stores all information of a call
type request struct {
	header     *codec.Header
	argValue   reflect.Value
	replyValue reflect.Value
}

func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{header: h}
	// TODO: now we don't know the type of request argValue
	// day 1, just suppose it's string
	req.argValue = reflect.New(reflect.TypeOf(""))
	if err = cc.ReadBody(req.argValue.Interface()); err != nil {
		log.Println("rpc server: read argValue err:", err)
	}
	return req, nil
}

func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
		return nil, err
	}
	return &h, nil
}

func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	// TODO, should call registered rpc methods to get the right replyValue
	// day 1, just print argValue and send a hello message
	defer wg.Done()
	log.Println(req.header, req.argValue.Elem())
	req.replyValue = reflect.ValueOf(fmt.Sprintf("geerpc resp %d", req.header.Seq))
	server.sendResponse(cc, req.header, req.replyValue.Interface(), sending)
}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
}
