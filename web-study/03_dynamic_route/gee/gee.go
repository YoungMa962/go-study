package gee

import (
	"net/http"
)

// HandlerFunc request handler
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	router *router
}

func NewEngine() *Engine {
	return &Engine{router: newRouter()}
}

//
// GET
// @Description:	GET请求
// @receiver engine
// @param pattern	请求路径
// @param handler	request handler
//
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

//
// POST
// @Description:	POST请求
// @receiver engine
// @param pattern	请求路径
// @param handler	request handler
//
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

//
// Run
// @Description:	开启一个http server并监听
// @receiver engine
// @param addr 		监听端口号
// @return err
//
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
