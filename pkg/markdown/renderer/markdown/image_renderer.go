package markdown

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type ImageRenderer struct {
}

// Render implements NodeRenderer.
func (*ImageRenderer) Render(r *Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	image, ok := node.(*ast.Image)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.Image, got '%T'", node)
	}

	if entering {
		r.w.AddIndentOnFirstWrite([]byte("!["))
		return ast.WalkContinue, nil
	}

	_, _ = fmt.Fprintf(r.w, "](%s", image.Destination)
	if len(image.Title) > 0 {
		_, _ = fmt.Fprintf(r.w, ` "%s"`, image.Title)
	}
	_, _ = r.w.Write([]byte{')'})

	return ast.WalkContinue, nil
}

var _ NodeRenderer = &ImageRenderer{}
