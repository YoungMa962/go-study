package gee

import "strings"

type node struct {
	pattern  string  // 待匹配路由
	part     string  // 路由中的一部分
	children []*node // 子节点
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

func newNode(part string, isWild bool) *node {
	return &node{part: part, isWild: isWild}
}

//
//  insert
//  @Description:递归查找每一层的节点，如果没有匹配到当前part的节点，则新建一个
//  @receiver n
//  @param pattern
//  @param parts
//  @param height
//
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = newNode(part, part[0] == ':' || part[0] == '*')
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

//
//  search
//  @Description:
//  @receiver n
//  @param parts
//  @param height
//  @return *node
//
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
