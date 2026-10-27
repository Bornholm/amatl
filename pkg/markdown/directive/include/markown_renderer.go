package include

import (
	"bytes"

	"github.com/Bornholm/amatl/pkg/markdown/directive"
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type MarkdownRenderer struct {
}

// Render implements markdown.NodeRenderer.
func (mr *MarkdownRenderer) Render(r *markdown.Render, directive *directive.Node, entering bool) (ast.WalkStatus, error) {
	includedSource, includedNode, err := includedContent(directive)
	if err != nil {
		return ast.WalkStop, errors.WithStack(err)
	}

	var buff bytes.Buffer

	if err := r.Renderer().Render(&buff, includedSource, includedNode); err != nil {
		return ast.WalkStop, errors.Wrap(err, "could not convert included source")
	}

	_, _ = r.Writer().Write(markdown.NewLineChar)
	_, _ = r.Writer().Write(markdown.NewLineChar)

	if _, err := r.Writer().Write(buff.Bytes()); err != nil {
		return ast.WalkStop, errors.WithStack(err)
	}

	return ast.WalkContinue, nil
}

var _ directive.MarkdownDirectiveRenderer = &MarkdownRenderer{}
