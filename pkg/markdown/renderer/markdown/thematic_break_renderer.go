package markdown

import (
	"github.com/yuin/goldmark/ast"
)

type ThematicBreakRenderer struct {
}

// Render implements NodeRenderer.
func (*ThematicBreakRenderer) Render(r *Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	_, _ = r.w.Write(ThematicBreakChars)

	return ast.WalkSkipChildren, nil
}

var _ NodeRenderer = &ThematicBreakRenderer{}
