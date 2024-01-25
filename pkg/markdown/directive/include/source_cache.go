package include

import "github.com/yuin/goldmark/ast"

type SourceCache struct {
	sources map[string][]byte
	nodes   map[string]ast.Node
}

func NewSourceCache() *SourceCache {
	return &SourceCache{
		sources: make(map[string][]byte),
		nodes:   make(map[string]ast.Node),
	}
}

func (c *SourceCache) Set(path string, data []byte, node ast.Node) {
	c.sources[path] = data
	c.nodes[path] = node
}

func (c *SourceCache) Get(path string) ([]byte, ast.Node, bool) {
	data, sourceExists := c.sources[path]
	node, nodeExists := c.nodes[path]

	if !sourceExists || !nodeExists {
		return nil, nil, false
	}

	return data, node, true
}
