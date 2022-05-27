package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

//
//  addRoute
//  @Description:  路由注册
//  @receiver r
//  @param method  请求类型
//  @param pattern 路由
//  @param handler 处理方法
//
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
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
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
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

	return nil, nil
}

//
//  parsePattern
//  @Description:
//  @param pattern
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
