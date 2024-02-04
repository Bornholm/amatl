package node

import (
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type DocumentRenderer struct {
}

// Render implements NodeRenderer.
func (*DocumentRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		return ast.WalkContinue, nil
	}

	if _, err := r.Writer().Write(markdown.NewLineChar); err != nil {
		return ast.WalkStop, errors.WithStack(err)
	}

	return ast.WalkContinue, nil
}

var _ markdown.NodeRenderer = &DocumentRenderer{}
