package node

import (
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type CodeSpanRenderer struct {
}

// Render implements NodeRenderer.
func (*CodeSpanRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	_, ok := node.(*ast.CodeSpan)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.CodeSpan, got '%T'", node)
	}

	if entering {
		_, _ = r.Writer().Write([]byte{'`'})
		return ast.WalkContinue, nil
	}

	_, _ = r.Writer().Write([]byte{'`'})
	return ast.WalkContinue, nil
}

var _ markdown.NodeRenderer = &CodeSpanRenderer{}
