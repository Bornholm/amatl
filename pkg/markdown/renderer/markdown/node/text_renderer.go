package node

import (
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type TextRenderer struct {
}

// Render implements NodeRenderer.
func (*TextRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	text, ok := node.(*ast.Text)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.Text, got '%T'", node)
	}

	if entering {
		text := text.Segment.Value(r.Source())
		_ = writeClean(r.Writer(), text)
		return ast.WalkContinue, nil
	}

	if text.SoftLineBreak() {
		char := markdown.SpaceChar
		if r.Renderer().SoftWraps() {
			char = markdown.NewLineChar
		}
		_, _ = r.Writer().Write(char)
	}

	if text.HardLineBreak() {
		if text.SoftLineBreak() {
			_, _ = r.Writer().Write(markdown.SpaceChar)
		}
		_, _ = r.Writer().Write(markdown.NewLineChar)
	}

	return ast.WalkContinue, nil
}

var _ markdown.NodeRenderer = &TextRenderer{}
