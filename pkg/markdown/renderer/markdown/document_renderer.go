package markdown

import (
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type DocumentRenderer struct {
}

// Render implements NodeRenderer.
func (*DocumentRenderer) Render(r *Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		return ast.WalkContinue, nil
	}

	if _, err := r.w.Write(NewLineChar); err != nil {
		return ast.WalkStop, errors.WithStack(err)
	}

	return ast.WalkContinue, nil
}

var _ NodeRenderer = &DocumentRenderer{}
