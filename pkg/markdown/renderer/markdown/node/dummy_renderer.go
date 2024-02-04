package node

import (
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/yuin/goldmark/ast"
)

type DummyRenderer struct {
}

// Render implements NodeRenderer.
func (*DummyRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

var _ markdown.NodeRenderer = &DummyRenderer{}
