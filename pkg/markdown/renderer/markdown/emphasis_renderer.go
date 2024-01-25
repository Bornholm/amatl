package markdown

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type EmphasisRenderer struct {
}

// Render implements NodeRenderer.
func (*EmphasisRenderer) Render(r *Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	emphasis, ok := node.(*ast.Emphasis)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.Emphasis, got '%T'", node)
	}

	var emWrapper []byte
	switch emphasis.Level {
	case 1:
		emWrapper = r.emphToken
	case 2:
		emWrapper = r.strongToken
	default:
		emWrapper = bytes.Repeat(r.emphToken, emphasis.Level)
	}

	return r.WrapNonEmptyContentWith(emWrapper, entering), nil
}

var _ NodeRenderer = &EmphasisRenderer{}
