package node

import (
	"bytes"

	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type EmphasisRenderer struct {
}

// Render implements NodeRenderer.
func (*EmphasisRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	emphasis, ok := node.(*ast.Emphasis)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.Emphasis, got '%T'", node)
	}

	var emWrapper []byte
	switch emphasis.Level {
	case 1:
		emWrapper = r.EmphToken()
	case 2:
		emWrapper = r.StrongToken()
	default:
		emWrapper = bytes.Repeat(r.EmphToken(), emphasis.Level)
	}

	return r.WrapNonEmptyContentWith(emWrapper, entering), nil
}

var _ markdown.NodeRenderer = &EmphasisRenderer{}
