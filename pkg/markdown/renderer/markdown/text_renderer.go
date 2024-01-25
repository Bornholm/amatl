package markdown

import (
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type TextRenderer struct {
}

// Render implements NodeRenderer.
func (*TextRenderer) Render(r *Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	text, ok := node.(*ast.Text)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.Text, got '%T'", node)
	}

	if entering {
		text := text.Segment.Value(r.source)
		_ = writeClean(r.w, text)
		return ast.WalkContinue, nil
	}

	if text.SoftLineBreak() {
		char := SpaceChar
		if r.mr.softWraps {
			char = NewLineChar
		}
		_, _ = r.w.Write(char)
	}

	if text.HardLineBreak() {
		if text.SoftLineBreak() {
			_, _ = r.w.Write(SpaceChar)
		}
		_, _ = r.w.Write(NewLineChar)
	}

	return ast.WalkContinue, nil
}

var _ NodeRenderer = &TextRenderer{}
