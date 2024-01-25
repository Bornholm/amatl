package markdown

import (
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type CodeSpanRenderer struct {
}

// Render implements NodeRenderer.
func (*CodeSpanRenderer) Render(r *Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	_, ok := node.(*ast.CodeSpan)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.CodeSpan, got '%T'", node)
	}

	if entering {
		_, _ = r.w.Write([]byte{'`'})
		return ast.WalkContinue, nil
	}

	_, _ = r.w.Write([]byte{'`'})
	return ast.WalkContinue, nil
}

var _ NodeRenderer = &CodeSpanRenderer{}
