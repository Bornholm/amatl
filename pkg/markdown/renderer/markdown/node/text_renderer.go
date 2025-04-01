package node

import (
	"slices"
	"strings"

	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
	extAST "github.com/yuin/goldmark/extension/ast"
)

type TextRenderer struct {
}

// Render implements NodeRenderer.
func (*TextRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	text, ok := node.(*ast.Text)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.Text, got '%T'", node)
	}

	inCodeSpan := hasAncestor[*ast.CodeSpan](node)
	inTable := hasAncestor[*extAST.TableCell](node)

	if entering {
		text := text.Segment.Value(r.Source())
		if inCodeSpan && inTable {
			text = escapeChars(text, '|')
		}
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

func escapeChars(text []byte, chars ...byte) []byte {
	var sb strings.Builder
	for _, r := range text {
		if slices.Contains(chars, r) {
			sb.WriteString(`\`)
		}

		sb.WriteByte(r)
	}
	return []byte(sb.String())
}

func hasAncestor[T any](n ast.Node) bool {
	parent := n.Parent()
	if parent == nil {
		return false
	}

	if _, ok := parent.(T); ok {
		return true
	}

	return hasAncestor[T](parent)
}
