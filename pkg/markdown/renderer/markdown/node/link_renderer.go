package node

import (
	"fmt"

	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type LinkRenderer struct {
}

// Render implements NodeRenderer.
func (*LinkRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	link, ok := node.(*ast.Link)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.Link, got '%T'", node)
	}

	if entering {
		r.Writer().AddIndentOnFirstWrite([]byte("["))
		return ast.WalkContinue, nil
	}

	_, _ = fmt.Fprintf(r.Writer(), "](%s", link.Destination)
	if len(link.Title) > 0 {
		_, _ = fmt.Fprintf(r.Writer(), ` "%s"`, link.Title)
	}
	_, _ = r.Writer().Write([]byte{')'})

	return ast.WalkContinue, nil
}

var _ markdown.NodeRenderer = &LinkRenderer{}
