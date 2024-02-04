package node

import (
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type AutoLinkRenderer struct {
}

// Render implements NodeRenderer.
func (*AutoLinkRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	link, ok := node.(*ast.AutoLink)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.AutoLink, got '%T'", node)
	}

	if entering {
		_, _ = r.Writer().Write(link.Label(r.Source()))
	}

	return ast.WalkContinue, nil
}

var _ markdown.NodeRenderer = &AutoLinkRenderer{}
