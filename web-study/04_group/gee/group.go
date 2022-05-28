package gee

import "log"

type RouterGroup struct {
	prefix     string
	parent     *RouterGroup  // 支持嵌套分组
	middleware []HandlerFunc // 中间件
	engine     Engine        // 提供router能力,所有RouterGroup共享同一个engine

}

func NewRouterGroup(prefix string, parent *RouterGroup, middleware []HandlerFunc, engine Engine) *RouterGroup {
	return &RouterGroup{prefix: prefix, parent: parent, middleware: middleware, engine: engine}
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := NewRouterGroup(group.prefix+prefix, group, nil, engine)
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

//
// GET
// @Description:	GET请求
// @receiver engine
// @param pattern	请求路径
// @param handler	request handler
//
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

//
// POST
// @Description:	POST请求
// @receiver engine
// @param pattern	请求路径
// @param handler	request handler
//
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Registy Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}
