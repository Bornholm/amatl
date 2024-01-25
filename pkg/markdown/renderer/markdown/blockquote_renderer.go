package markdown

import (
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type BlockquoteRenderer struct {
}

// Render implements NodeRenderer.
func (*BlockquoteRenderer) Render(r *Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	_, ok := node.(*ast.Blockquote)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.Blockquote, got '%T'", node)
	}

	if entering {
		r.w.PushIndent(BlockquoteChars)
		if node.Parent() != nil && node.Parent().Kind() == ast.KindListItem &&
			node.PreviousSibling() == nil {
			_, _ = r.w.Write(BlockquoteChars)
		}
	} else {
		r.w.PopIndent()
	}

	return ast.WalkContinue, nil
}

var _ NodeRenderer = &BlockquoteRenderer{}
