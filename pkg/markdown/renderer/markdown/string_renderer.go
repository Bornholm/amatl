package markdown

import (
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type StringRenderer struct {
}

// Render implements NodeRenderer.
func (*StringRenderer) Render(r *Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	str, ok := node.(*ast.String)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.String, got '%T'", node)
	}

	if entering {
		_, _ = r.w.Write(str.Value)
	}

	return ast.WalkContinue, nil
}

var _ NodeRenderer = &StringRenderer{}
