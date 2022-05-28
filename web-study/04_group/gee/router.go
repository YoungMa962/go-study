package gee

import (
	"net/http"
	"strings"
)

//
//  router
//  @Description: 路由实现，核心方法 1、路由注册 2、路由匹配
//
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),       // 前缀树用于路由匹配
		handlers: make(map[string]HandlerFunc), // 路由地址和handler的映射map
	}
}

//
//  addRoute
//  @Description:  路由注册
// 					1、在前缀树中添加相应的路由节点
// 					2、注册对应的路由地址的handler
//  @receiver r
//  @param method  请求类型
//  @param pattern 路由
//  @param handler 处理方法
//
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	// 在前缀树中添加相应的路由节点
	r.roots[method].insert(pattern, parts, 0)
	// 注册对应的路由地址的handler
	r.handlers[key] = handler
}

//
//  getRoute
//  @Description: 路由匹配
//  @receiver r
//  @param method 请求类型
//  @param path	  请求路径
//  @return *node
//  @return map[string]string
//
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	if root, ok := r.roots[method]; ok {
		if n := root.search(searchParts, 0); n != nil {
			parts := parsePattern(n.pattern)
			for index, part := range parts {
				if part[0] == ':' {
					params[part[1:]] = searchParts[index]
				}
				if part[0] == '*' && len(part) > 1 {
					params[part[1:]] = strings.Join(searchParts[index:], "/")
					break
				}
			}
			return n, params
		}
	}
	return nil, nil
}

//
//  parsePattern
//  @Description:	解析路由地址 生成路由数组
//  @param pattern	路由地址
//  @return []string
//
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method(), c.Path())
	if n != nil {
		c.SetParams(params)
		key := c.Method() + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.StringResponse(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
