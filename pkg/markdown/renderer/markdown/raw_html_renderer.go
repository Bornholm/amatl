package markdown

import (
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type RawHTMLRenderer struct {
}

// Render implements NodeRenderer.
func (*RawHTMLRenderer) Render(r *Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	rawHTML, ok := node.(*ast.RawHTML)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.RawHTML, got '%T'", node)
	}

	if !entering {
		return ast.WalkContinue, nil
	}

	for i := 0; i < rawHTML.Segments.Len(); i++ {
		segment := rawHTML.Segments.At(i)
		_, _ = r.w.Write(segment.Value(r.source))
	}

	return ast.WalkSkipChildren, nil
}

var _ NodeRenderer = &RawHTMLRenderer{}
