package markdown

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type HTMLBlockRenderer struct {
}

// Render implements NodeRenderer.
func (*HTMLBlockRenderer) Render(r *Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	html, ok := node.(*ast.HTMLBlock)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.HTMLBlock, got '%T'", node)
	}

	if !entering {
		return ast.WalkContinue, nil
	}

	var segments []text.Segment
	for i := 0; i < node.Lines().Len(); i++ {
		segments = append(segments, node.Lines().At(i))
	}

	if html.ClosureLine.Len() != 0 {
		segments = append(segments, html.ClosureLine)
	}
	for i, s := range segments {
		o := s.Value(r.source)
		if i == len(segments)-1 {
			o = bytes.TrimSuffix(o, []byte("\n"))
		}
		_, _ = r.w.Write(o)
	}
	return ast.WalkSkipChildren, nil

}

var _ NodeRenderer = &HTMLBlockRenderer{}
