package node

import (
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type RawHTMLRenderer struct {
}

// Render implements NodeRenderer.
func (*RawHTMLRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	rawHTML, ok := node.(*ast.RawHTML)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.RawHTML, got '%T'", node)
	}

	if !entering {
		return ast.WalkContinue, nil
	}

	for i := 0; i < rawHTML.Segments.Len(); i++ {
		segment := rawHTML.Segments.At(i)
		_, _ = r.Writer().Write(segment.Value(r.Source()))
	}

	return ast.WalkSkipChildren, nil
}

var _ markdown.NodeRenderer = &RawHTMLRenderer{}
