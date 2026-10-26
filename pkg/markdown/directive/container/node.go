package container

import (
	"github.com/yuin/goldmark/ast"
)

// Kind is the node kind of container directives
var Kind = ast.NewNodeKind("ContainerDirective")

// Node is a block-level directive wrapping other blocks:
//
//	:::name{#id .class key=value}
//	Some **markdown** content
//	:::
type Node struct {
	ast.BaseBlock
	name        string
	fenceLength int
	closed      bool
}

// Name returns the directive name
func (n *Node) Name() string {
	return n.name
}

// FenceLength returns the number of colons used to open the container
func (n *Node) FenceLength() int {
	return n.fenceLength
}

// Kind implements ast.Node.
func (n *Node) Kind() ast.NodeKind {
	return Kind
}

// Dump implements ast.Node.
func (n *Node) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{
		"Name": n.name,
	}, nil)
}

func NewNode(name string, fenceLength int) *Node {
	return &Node{
		name:        name,
		fenceLength: fenceLength,
	}
}

var _ ast.Node = &Node{}
