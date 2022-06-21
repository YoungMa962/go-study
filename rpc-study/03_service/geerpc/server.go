package geerpc

import (
	"codec"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
)

type Server struct {
	// 服务列表
	serviceMap sync.Map
}

func NewServer() *Server {
	return &Server{}
}

var defaultServer = NewServer()

// Accept 接收clint 请求
func Accept(lis net.Listener) { defaultServer.accept(lis) }

func (server *Server) accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
			return
		}
		go server.serveConn(conn)
	}
}

func (server *Server) serveConn(conn io.ReadWriteCloser) {
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

// Registry publishes in the server the set of methods of the
func Registry(receiver interface{}) error {
	return defaultServer.register(receiver)
}
func (server *Server) register(receiver interface{}) error {
	s := newService(receiver)
	if _, dup := server.serviceMap.LoadOrStore(s.name, s); dup {
		return errors.New("rpc: service already defined: " + s.name)
	}
	return nil
}

func (server *Server) findService(serviceMethod string) (mService *service, mType *methodType, err error) {
	split := strings.Split(serviceMethod, ".")
	serviceName, methodName := split[0], split[1]
	ser, ok := server.serviceMap.Load(serviceName)
	if !ok {
		log.Fatalf("404 Service{%s}Not Find", serviceName)
	}
	mService = ser.(*service)
	mType = mService.method[methodName]
	if mType == nil {
		err = errors.New("rpc server: can't find method " + methodName)
	}
	return
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
	mtype      *methodType
	svc        *service
}

func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{header: h}
	req.svc, req.mtype, err = server.findService(h.ServiceMethod)
	if err != nil {
		return req, err
	}

	req.argValue = req.mtype.newArgValue()
	req.replyValue = req.mtype.newReplyValue()

	argvi := req.argValue.Interface()
	if req.argValue.Type().Kind() != reflect.Ptr {
		argvi = req.argValue.Addr().Interface()
	}
	if err = cc.ReadBody(argvi); err != nil {
		log.Println("rpc server: read body err:", err)
		return req, err
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
	defer wg.Done()
	err := req.svc.call(req.mtype, req.argValue, req.replyValue)
	if err != nil {
		req.header.ResError = err.Error()
		server.sendResponse(cc, req.header, invalidRequest, sending)
		return
	}
	server.sendResponse(cc, req.header, req.replyValue.Interface(), sending)
}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
}
