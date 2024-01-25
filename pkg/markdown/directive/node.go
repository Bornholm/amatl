package directive

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type Type string

type Node struct {
	ast.BaseInline
	directiveType Type

	value *ast.Text
}

var KindDirective = ast.NewNodeKind("Directive")

// Dump implements ast.Node.
func (n *Node) Dump(source []byte, level int) {
	segment := n.value.Segment
	m := map[string]string{
		"Value": string(segment.Value(source)),
	}
	ast.DumpHelper(n, source, level, m, nil)
}

// Kind implements ast.Node.
func (n *Node) Kind() ast.NodeKind {
	return KindDirective
}

func (n *Node) DirectiveType() Type {
	return n.directiveType
}

func parseDirective(raw []byte, value *ast.Text) *Node {
	var (
		pos           int
		start         int
		state         int
		rawName       []byte
		rawAttributes []byte
	)
LOOP:
	for pos = range raw {
		switch state {
		case 0: // Parse directive name
			if pos == 0 && raw[pos] != ':' {
				break LOOP
			}

			if raw[pos] != '{' {
				continue
			}

			rawName = raw[1:pos]
			state = 1
			start = pos

		case 1: // Parse attributes
			if pos == start && raw[pos] != '{' {
				break LOOP
			}

			if raw[pos] == '}' {
				rawAttributes = raw[start : pos+1]
				break LOOP
			}
		}

	}
	if rawName == nil {
		return nil
	}

	var (
		attributes parser.Attributes
		ok         bool
	)

	if rawAttributes != nil {
		reader := text.NewReader(rawAttributes)
		attributes, ok = parser.ParseAttributes(reader)
		if !ok {
			return nil
		}
	}

	node := &Node{
		directiveType: Type(rawName),
		BaseInline:    ast.BaseInline{},
		value:         value,
	}

	for _, attr := range attributes {
		value := attr.Value.([]byte)
		node.SetAttribute(attr.Name, string(value))
	}

	return node
}

var _ ast.Node = &Node{}
