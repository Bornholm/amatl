package markdown

import (
	"github.com/yuin/goldmark/ast"
)

type StrikethroughRenderer struct {
}

// Render implements NodeRenderer.
func (*StrikethroughRenderer) Render(r *Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return r.WrapNonEmptyContentWith(StrikeThroughChars, entering), nil
}

var _ NodeRenderer = &StrikethroughRenderer{}
