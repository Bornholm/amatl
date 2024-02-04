package node

import (
	"bytes"

	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/yuin/goldmark/ast"
)

type CodeBlockRenderer struct {
}

// Render implements NodeRenderer.
func (*CodeBlockRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	_, _ = r.Writer().Write(markdown.CodeBlockChars)

	var lang []byte
	if fencedNode, isFenced := node.(*ast.FencedCodeBlock); isFenced && fencedNode.Info != nil {
		lang = fencedNode.Info.Text(r.Source())
		_, _ = r.Writer().Write(lang)
		for _, elt := range bytes.Fields(lang) {
			elt = bytes.TrimSpace(bytes.TrimLeft(elt, ". "))
			if len(elt) == 0 {
				continue
			}
			lang = elt
			break
		}
	}

	_, _ = r.Writer().Write(markdown.NewLineChar)
	codeBuf := bytes.Buffer{}
	for i := 0; i < node.Lines().Len(); i++ {
		line := node.Lines().At(i)
		_, _ = codeBuf.Write(line.Value(r.Source()))
	}

	if formatCode, ok := r.Renderer().Formatters()[noAllocString(lang)]; ok {
		code := formatCode(codeBuf.Bytes())
		if !bytes.HasSuffix(code, markdown.NewLineChar) {
			// Ensure code sample ends with a newline.
			code = append(code, markdown.NewLineChar...)
		}
		_, _ = r.Writer().Write(code)
	} else {
		_, _ = r.Writer().Write(codeBuf.Bytes())
	}

	_, _ = r.Writer().Write(markdown.CodeBlockChars)

	return ast.WalkSkipChildren, nil
}

var _ markdown.NodeRenderer = &CodeBlockRenderer{}
