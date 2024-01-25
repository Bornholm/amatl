package markdown

import (
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type AutoLinkRenderer struct {
}

// Render implements NodeRenderer.
func (*AutoLinkRenderer) Render(r *Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	link, ok := node.(*ast.AutoLink)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.AutoLink, got '%T'", node)
	}

	if entering {
		_, _ = r.w.Write(link.Label(r.source))
	}

	return ast.WalkContinue, nil
}

var _ NodeRenderer = &AutoLinkRenderer{}
