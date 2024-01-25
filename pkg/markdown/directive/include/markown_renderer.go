package include

import (
	"bytes"

	"forge.cadoles.com/wpetit/amatl/pkg/markdown/directive"
	"forge.cadoles.com/wpetit/amatl/pkg/markdown/renderer/markdown"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type MarkdownRenderer struct {
	Cache *SourceCache
}

// Render implements markdown.NodeRenderer.
func (mr *MarkdownRenderer) Render(r *markdown.Render, directive *directive.Node, entering bool) (ast.WalkStatus, error) {
	rawPath, exists := directive.AttributeString("path")
	if !exists {
		return ast.WalkStop, errors.Errorf("could not find 'path' attribute in directive '%s'", directive.Kind())
	}

	path, ok := rawPath.(string)
	if !ok {
		panic(errors.Errorf("unexpected type '%T' for 'path' attribute", rawPath))
	}

	includedSource, includedNode, exists := mr.Cache.Get(path)
	if !exists {
		panic(errors.Errorf("could not find source associated with path '%s'", path))
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
