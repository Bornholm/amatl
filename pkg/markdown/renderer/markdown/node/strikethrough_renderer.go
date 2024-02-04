package node

import (
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/yuin/goldmark/ast"
)

type StrikethroughRenderer struct {
}

// Render implements NodeRenderer.
func (*StrikethroughRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return r.WrapNonEmptyContentWith(markdown.StrikeThroughChars, entering), nil
}

var _ markdown.NodeRenderer = &StrikethroughRenderer{}
