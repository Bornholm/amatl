package node

import (
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/yuin/goldmark/ast"
)

type ThematicBreakRenderer struct {
}

// Render implements NodeRenderer.
func (*ThematicBreakRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	_, _ = r.Writer().Write(markdown.ThematicBreakChars)

	return ast.WalkSkipChildren, nil
}

var _ markdown.NodeRenderer = &ThematicBreakRenderer{}
