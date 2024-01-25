package include

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Node struct {
	wrapped ast.Node
	source  []byte
}

// AppendChild implements ast.Node.
func (n *Node) AppendChild(self ast.Node, child ast.Node) {
	n.wrapped.AppendChild(self, child)
}

// Attribute implements ast.Node.
func (n *Node) Attribute(name []byte) (interface{}, bool) {
	return n.wrapped.Attribute(name)
}

// AttributeString implements ast.Node.
func (n *Node) AttributeString(name string) (interface{}, bool) {
	return n.wrapped.AttributeString(name)
}

// Attributes implements ast.Node.
func (n *Node) Attributes() []ast.Attribute {
	return n.wrapped.Attributes()
}

// ChildCount implements ast.Node.
func (n *Node) ChildCount() int {
	return n.wrapped.ChildCount()
}

// Dump implements ast.Node.
func (n *Node) Dump(source []byte, level int) {
	n.wrapped.Dump(n.source, level)
}

// FirstChild implements ast.Node.
func (n *Node) FirstChild() ast.Node {
	if n.wrapped == nil {
		return nil
	}

	return &Node{
		wrapped: n.wrapped.FirstChild(),
		source:  n.source,
	}
}

// HasBlankPreviousLines implements ast.Node.
func (n *Node) HasBlankPreviousLines() bool {
	return n.wrapped.HasBlankPreviousLines()
}

// HasChildren implements ast.Node.
func (n *Node) HasChildren() bool {
	return n.wrapped.HasChildren()
}

// InsertAfter implements ast.Node.
func (n *Node) InsertAfter(self ast.Node, v1 ast.Node, insertee ast.Node) {
	n.wrapped.InsertAfter(self, v1, insertee)
}

// InsertBefore implements ast.Node.
func (n *Node) InsertBefore(self ast.Node, v1 ast.Node, insertee ast.Node) {
	n.wrapped.InsertBefore(self, v1, insertee)
}

// IsRaw implements ast.Node.
func (n *Node) IsRaw() bool {
	return n.wrapped.IsRaw()
}

// Kind implements ast.Node.
func (n *Node) Kind() ast.NodeKind {
	return n.wrapped.Kind()
}

// LastChild implements ast.Node.
func (n *Node) LastChild() ast.Node {
	if n.wrapped == nil {
		return nil
	}

	return &Node{
		wrapped: n.wrapped.LastChild(),
		source:  n.source,
	}
}

// Lines implements ast.Node.
func (n *Node) Lines() *text.Segments {
	return n.wrapped.Lines()
}

// NextSibling implements ast.Node.
func (n *Node) NextSibling() ast.Node {
	if n.wrapped == nil {
		return nil
	}

	if n.wrapped.Kind() == ast.KindHeading {
		return n.wrapped
	}

	return &Node{
		wrapped: n.wrapped.NextSibling(),
		source:  n.source,
	}
}

// OwnerDocument implements ast.Node.
func (n *Node) OwnerDocument() *ast.Document {
	return n.wrapped.OwnerDocument()
}

// Parent implements ast.Node.
func (n *Node) Parent() ast.Node {
	if n.wrapped == nil {
		return nil
	}

	if n.wrapped.Kind() == ast.KindHeading {
		return n.wrapped
	}

	return &Node{
		wrapped: n.wrapped.Parent(),
		source:  n.source,
	}
}

// PreviousSibling implements ast.Node.
func (n *Node) PreviousSibling() ast.Node {
	if n.wrapped == nil {
		return nil
	}

	return &Node{
		wrapped: n.wrapped.PreviousSibling(),
		source:  n.source,
	}
}

// RemoveAttributes implements ast.Node.
func (n *Node) RemoveAttributes() {
	if n.wrapped == nil {
		return
	}

	n.wrapped.RemoveAttributes()
}

// RemoveChild implements ast.Node.
func (n *Node) RemoveChild(self ast.Node, child ast.Node) {
	if n.wrapped == nil {
		return
	}

	n.wrapped.RemoveChild(self, child)
}

// RemoveChildren implements ast.Node.
func (n *Node) RemoveChildren(self ast.Node) {
	if n.wrapped == nil {
		return
	}

	n.wrapped.RemoveChildren(self)
}

// ReplaceChild implements ast.Node.
func (n *Node) ReplaceChild(self ast.Node, v1 ast.Node, insertee ast.Node) {
	if n.wrapped == nil {
		return
	}

	n.wrapped.ReplaceChild(self, v1, insertee)
}

// SetAttribute implements ast.Node.
func (n *Node) SetAttribute(name []byte, value interface{}) {
	n.wrapped.SetAttribute(name, value)
}

// SetAttributeString implements ast.Node.
func (n *Node) SetAttributeString(name string, value interface{}) {
	n.wrapped.SetAttributeString(name, value)
}

// SetBlankPreviousLines implements ast.Node.
func (n *Node) SetBlankPreviousLines(v bool) {
	n.wrapped.SetBlankPreviousLines(v)
}

// SetLines implements ast.Node.
func (n *Node) SetLines(segments *text.Segments) {
	n.wrapped.SetLines(segments)
}

// SetNextSibling implements ast.Node.
func (n *Node) SetNextSibling(next ast.Node) {
	if n.wrapped == nil {
		return
	}

	n.wrapped.SetNextSibling(next)
}

// SetParent implements ast.Node.
func (n *Node) SetParent(parent ast.Node) {
	if n.wrapped == nil {
		return
	}

	n.wrapped.SetParent(parent)
}

// SetPreviousSibling implements ast.Node.
func (n *Node) SetPreviousSibling(ast.Node) {
	if n.wrapped == nil {
		return
	}

	n.wrapped.SetPreviousSibling(n)
}

// SortChildren implements ast.Node.
func (n *Node) SortChildren(comparator func(n1 ast.Node, n2 ast.Node) int) {
	n.wrapped.SortChildren(comparator)
}

// Text implements ast.Node.
func (n *Node) Text(source []byte) []byte {
	return n.wrapped.Text(n.source)
}

// Type implements ast.Node.
func (n *Node) Type() ast.NodeType {
	return n.wrapped.Type()
}

var _ ast.Node = &Node{}
