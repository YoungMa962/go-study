package gee

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	req        *http.Request
	writer     http.ResponseWriter
	path       string            // 请求路径
	method     string            // 请求类型
	params     map[string]string // 请求参数
	statusCode int               // 响应值
}

// 构造函数 封装 writer 和 req
func newContext(writer http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		writer: writer,
		req:    req,
		path:   req.URL.Path,
		method: req.Method,
	}
}

// SetStatusCode 设置 status code
func (c *Context) SetStatusCode(statusCode int) {
	c.statusCode = statusCode
	c.writer.WriteHeader(statusCode)
}

// SetHeader 设置请求头
func (c *Context) SetHeader(key, value string) {
	c.writer.Header().Set(key, value)
}

func (c *Context) Body() map[string]string {
	jsonString, _ := ioutil.ReadAll(c.req.Body)
	m := make(map[string]string)
	err := json.Unmarshal(jsonString, &m)
	fmt.Printf("request body : %v\n", m)
	if err != nil {
		return nil
	} else {
		return m
	}
}

// Form get Form PARAMETERS
func (c *Context) Form(key string) string {
	return c.req.FormValue(key)
}

// Query get QUERY PARAMETERS
func (c *Context) Query(key string) string {
	return c.req.URL.Query().Get(key)
}

// DataResponse byte data res
func (c *Context) DataResponse(code int, data []byte) {
	c.SetStatusCode(code)
	c.writer.Write(data)
}

// StringResponse text/plain
func (c *Context) StringResponse(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatusCode(code)
	c.writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSONResponse application/json
func (c *Context) JSONResponse(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatusCode(code)
	encoder := json.NewEncoder(c.writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.writer, err.Error(), 500)
	}
}

// HTMLResponse text/html
func (c *Context) HTMLResponse(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatusCode(code)
	c.writer.Write([]byte(html))
}

// Param
// @Description:
// @receiver c
// @param key
// @return string
//
func (c *Context) Param(key string) string {
	value, _ := c.params[key]
	return value
}

func (c *Context) Req() *http.Request {
	return c.req
}

func (c *Context) SetReq(req *http.Request) {
	c.req = req
}

func (c *Context) Writer() http.ResponseWriter {
	return c.writer
}

func (c *Context) SetWriter(writer http.ResponseWriter) {
	c.writer = writer
}

func (c *Context) Path() string {
	return c.path
}

func (c *Context) SetPath(path string) {
	c.path = path
}

func (c *Context) Method() string {
	return c.method
}

func (c *Context) SetMethod(method string) {
	c.method = method
}

func (c *Context) Params() map[string]string {
	return c.params
}

func (c *Context) SetParams(params map[string]string) {
	c.params = params
}
