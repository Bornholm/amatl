package markdown

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type LinkRenderer struct {
}

// Render implements NodeRenderer.
func (*LinkRenderer) Render(r *Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	link, ok := node.(*ast.Link)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.Link, got '%T'", node)
	}

	if entering {
		r.w.AddIndentOnFirstWrite([]byte("["))
		return ast.WalkContinue, nil
	}

	_, _ = fmt.Fprintf(r.w, "](%s", link.Destination)
	if len(link.Title) > 0 {
		_, _ = fmt.Fprintf(r.w, ` "%s"`, link.Title)
	}
	_, _ = r.w.Write([]byte{')'})

	return ast.WalkContinue, nil
}

var _ NodeRenderer = &LinkRenderer{}
