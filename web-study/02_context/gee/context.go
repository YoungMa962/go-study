package gee

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	writer     http.ResponseWriter
	req        *http.Request
	path       string
	method     string
	statusCode int
}

func (c *Context) Path() string {
	return c.path
}

func (c *Context) Method() string {
	return c.method
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
