package gee

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// HandlerFunc request handler
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	router *router        // 提供路由能力
	root   *RouterGroup   // 最顶层分组
	groups []*RouterGroup // 存储所有的groups
}

func NewEngine() *Engine {
	engine := &Engine{router: newRouter()}
	engine.root = NewRouterGroup("", nil, nil, *engine)
	engine.groups = []*RouterGroup{engine.root}
	return engine
}

func (engine *Engine) Use(middlewares ...HandlerFunc) {
	engine.root.Use(middlewares...)
}

func (engine *Engine) Group(prefix string) *RouterGroup {
	routerGroup := engine.root.Group(prefix)
	engine.groups = append(engine.groups, routerGroup)
	return routerGroup
}

//
// GET
// @Description:	GET请求
// @receiver engine
// @param pattern	请求路径
// @param handler	request handler
//
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.root.GET(pattern, handler)
}

//
// POST
// @Description:	POST请求
// @receiver engine
// @param pattern	请求路径
// @param handler	request handler
//
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.root.POST(pattern, handler)
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

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 添加 中间件函数到执行列表
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) && group.middlewares != nil {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}

func Logger() HandlerFunc {
	return func(context *Context) {
		// Start timer
		t := time.Now()
		// Process request
		context.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", context.StatusCode(), context.Req().RequestURI, time.Since(t))
	}
}
