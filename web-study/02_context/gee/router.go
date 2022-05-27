package gee

import (
	"log"
	"net/http"
)

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

// 添加路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.method + "-" + c.path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.StringResponse(http.StatusNotFound, "404 NOT FOUND: %s\n", c.path)
	}
}
