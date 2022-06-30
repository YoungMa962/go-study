package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

var _ http.Handler = (*Proxy)(nil) // 实现io.Closer 接口
type Proxy struct {
}

func (p *Proxy) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	// 判断协议类型
	if request.Method != "CONNECT" {
		p.HTTPHandler(writer, request)
	} else {
		// p.HTTPSHandler(writer, request)
	}
}

func (p *Proxy) HTTPHandler(writer http.ResponseWriter, request *http.Request) {
	printHttpRequest(request)
	transport := http.DefaultTransport
	outReq := new(http.Request) // 用于转发请求
	*outReq = *request          // 复制客户端请求

	// 发送请求
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		writer.WriteHeader(http.StatusBadGateway)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	// 返回信息给客户端
	// 返回状态码
	writer.WriteHeader(res.StatusCode)

	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(res.Body)
	backBody := buf.String()

	// 返回请求体
	_, _ = io.WriteString(writer, backBody)
	_ = res.Body.Close()
}

func (p *Proxy) HTTPSHandler(writer http.ResponseWriter, request *http.Request) {
	printHttpRequest(request)
	// 和目标服务器创建连接
	dst, err := net.DialTimeout("tcp", request.Host, 10*time.Second)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusServiceUnavailable)
		return
	}
	writer.WriteHeader(http.StatusOK)
	hijacker, ok := writer.(http.Hijacker)
	if !ok {
		http.Error(writer, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(dst, conn)
	go transfer(conn, dst)
}

func transfer(dst io.WriteCloser, src io.ReadCloser) {
	defer dst.Close()
	defer src.Close()
	_, _ = io.Copy(dst, src)
}
func printHttpRequest(request *http.Request) {
	fmt.Printf("%s\t%s\t%s\n", request.Method, request.URL.RequestURI(), request.Proto)
	fmt.Printf("Host:%v\n", request.Host)
	for k, v := range request.Header {
		fmt.Printf("%s:%s\n", k, v)
	}
	fmt.Println("===============================================")
}
