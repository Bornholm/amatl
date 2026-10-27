package include

import (
	"github.com/pkg/errors"
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

// includedContent returns the source and the node the include transformer
// attached to the given directive node
func includedContent(n ast.Node) ([]byte, ast.Node, error) {
	includedSource, sourceExists := IncludedSource(n)
	includedNode, nodeExists := IncludedNode(n)

	if !sourceExists || !nodeExists {
		rawURL, _ := getNodeURLAttribute(n)
		return nil, nil, errors.Errorf("could not find source associated with path '%s'", rawURL)
	}

	return includedSource, includedNode, nil
}
