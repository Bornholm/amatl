package mermaid

import (
	"bytes"

	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/yuin/goldmark/ast"
)

type BlockNodeRenderer struct {
}

// Render implements markdown.NodeRenderer.
func (*BlockNodeRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	_, _ = r.Writer().Write(markdown.NewLineChar)
	_, _ = r.Writer().Write(markdown.NewLineChar)
	_, _ = r.Writer().Write(markdown.CodeBlockChars)

	lang := []byte("mermaid")
	_, _ = r.Writer().Write(lang)
	for _, elt := range bytes.Fields(lang) {
		elt = bytes.TrimSpace(bytes.TrimLeft(elt, ". "))
		if len(elt) == 0 {
			continue
		}
		break
	}

	_, _ = r.Writer().Write(markdown.NewLineChar)
	codeBuf := bytes.Buffer{}
	for i := 0; i < node.Lines().Len(); i++ {
		line := node.Lines().At(i)
		_, _ = codeBuf.Write(line.Value(r.Source()))
	}

	_, _ = r.Writer().Write(codeBuf.Bytes())
	_, _ = r.Writer().Write(markdown.CodeBlockChars)

	return ast.WalkSkipChildren, nil
}

var _ markdown.NodeRenderer = &BlockNodeRenderer{}

type ScriptBlockNodeRenderer struct {
}

// Render implements markdown.NodeRenderer.
func (*ScriptBlockNodeRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

var _ markdown.NodeRenderer = &ScriptBlockNodeRenderer{}
