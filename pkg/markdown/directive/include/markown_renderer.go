package include

import (
	"bytes"

	"github.com/Bornholm/amatl/pkg/markdown/directive"
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type MarkdownRenderer struct {
	Cache *SourceCache
}

// Render implements markdown.NodeRenderer.
func (mr *MarkdownRenderer) Render(r *markdown.Render, directive *directive.Node, entering bool) (ast.WalkStatus, error) {
	rawURL, err := getNodeURLAttribute(directive)
	if err != nil {
		return ast.WalkStop, errors.WithStack(err)
	}

	includedSource, includedNode, exists := mr.Cache.Get(rawURL)
	if !exists {
		panic(errors.Errorf("could not find source associated with url '%s'", rawURL))
	}

	var buff bytes.Buffer

	if err := r.Renderer().Render(&buff, includedSource, includedNode); err != nil {
		panic(errors.Wrap(err, "could not convert included source"))
	}

	_, _ = r.Writer().Write(markdown.NewLineChar)
	_, _ = r.Writer().Write(markdown.NewLineChar)

	if _, err := r.Writer().Write(buff.Bytes()); err != nil {
		panic(errors.WithStack(err))
	}

	return ast.WalkContinue, nil
}

var _ directive.MarkdownDirectiveRenderer = &MarkdownRenderer{}
