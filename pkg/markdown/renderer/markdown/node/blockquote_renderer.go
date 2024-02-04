package node

import (
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type BlockquoteRenderer struct {
}

// Render implements NodeRenderer.
func (*BlockquoteRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	_, ok := node.(*ast.Blockquote)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.Blockquote, got '%T'", node)
	}

	if entering {
		r.Writer().PushIndent(markdown.BlockquoteChars)
		if node.Parent() != nil && node.Parent().Kind() == ast.KindListItem &&
			node.PreviousSibling() == nil {
			_, _ = r.Writer().Write(markdown.BlockquoteChars)
		}
	} else {
		r.Writer().PopIndent()
	}

	return ast.WalkContinue, nil
}

var _ markdown.NodeRenderer = &BlockquoteRenderer{}
