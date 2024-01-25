package markdown

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
)

type CodeBlockRenderer struct {
}

// Render implements NodeRenderer.
func (*CodeBlockRenderer) Render(r *Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	_, _ = r.w.Write(CodeBlockChars)

	var lang []byte
	if fencedNode, isFenced := node.(*ast.FencedCodeBlock); isFenced && fencedNode.Info != nil {
		lang = fencedNode.Info.Text(r.source)
		_, _ = r.w.Write(lang)
		for _, elt := range bytes.Fields(lang) {
			elt = bytes.TrimSpace(bytes.TrimLeft(elt, ". "))
			if len(elt) == 0 {
				continue
			}
			lang = elt
			break
		}
	}

	_, _ = r.w.Write(NewLineChar)
	codeBuf := bytes.Buffer{}
	for i := 0; i < node.Lines().Len(); i++ {
		line := node.Lines().At(i)
		_, _ = codeBuf.Write(line.Value(r.source))
	}

	if formatCode, ok := r.mr.formatters[noAllocString(lang)]; ok {
		code := formatCode(codeBuf.Bytes())
		if !bytes.HasSuffix(code, NewLineChar) {
			// Ensure code sample ends with a newline.
			code = append(code, NewLineChar...)
		}
		_, _ = r.w.Write(code)
	} else {
		_, _ = r.w.Write(codeBuf.Bytes())
	}

	_, _ = r.w.Write(CodeBlockChars)

	return ast.WalkSkipChildren, nil
}

var _ NodeRenderer = &CodeBlockRenderer{}
