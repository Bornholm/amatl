package markdown

import (
	"github.com/yuin/goldmark/ast"
)

type DummyRenderer struct {
}

// Render implements NodeRenderer.
func (*DummyRenderer) Render(r *Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

var _ NodeRenderer = &DummyRenderer{}
