package include

import (
	"github.com/yuin/goldmark/ast"
)

const (
	attrIncludedNode   = "includedNode"
	attrIncludedSource = "includedSource"
)

func setIncludedNode(n ast.Node, includedNode ast.Node) {
	n.SetAttributeString(attrIncludedNode, includedNode)
}

func IncludedNode(n ast.Node) (ast.Node, bool) {
	raw, exists := n.AttributeString(attrIncludedNode)
	if !exists {
		return nil, false
	}

	includedNode, ok := raw.(ast.Node)
	if !ok {
		return nil, false
	}

	return includedNode, true
}

func setIncludedSource(n ast.Node, includedSource []byte) {
	n.SetAttributeString(attrIncludedSource, includedSource)
}

func IncludedSource(n ast.Node) ([]byte, bool) {
	raw, exists := n.AttributeString(attrIncludedSource)
	if !exists {
		return nil, false
	}

	includedSource, ok := raw.([]byte)
	if !ok {
		return nil, false
	}

	return includedSource, true
}
