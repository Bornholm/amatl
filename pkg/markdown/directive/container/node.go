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
	// column of the opening fence, from the beginning of the line
	column int
	closed bool
	// source offsets of the opening fence start and of the closing fence end,
	// -1 when the container has no closing fence
	start int
	stop  int
}

// Name returns the directive name
func (n *Node) Name() string {
	return n.name
}

// FenceLength returns the number of colons used to open the container
func (n *Node) FenceLength() int {
	return n.fenceLength
}

// SourceRange returns the [start, stop) byte range of the container in the
// source, fences included. stop is -1 when no closing fence ends the
// container, for example at the end of the document.
func (n *Node) SourceRange() (start int, stop int) {
	return n.start, n.stop
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

func NewNode(name string, fenceLength int, column int) *Node {
	return &Node{
		name:        name,
		fenceLength: fenceLength,
		column:      column,
		start:       -1,
		stop:        -1,
	}
}

var _ ast.Node = &Node{}
