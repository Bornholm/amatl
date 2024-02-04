package node

import (
	"fmt"

	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type ImageRenderer struct {
}

// Render implements NodeRenderer.
func (*ImageRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	image, ok := node.(*ast.Image)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.Image, got '%T'", node)
	}

	if entering {
		r.Writer().AddIndentOnFirstWrite([]byte("!["))
		return ast.WalkContinue, nil
	}

	_, _ = fmt.Fprintf(r.Writer(), "](%s", image.Destination)
	if len(image.Title) > 0 {
		_, _ = fmt.Fprintf(r.Writer(), ` "%s"`, image.Title)
	}
	_, _ = r.Writer().Write([]byte{')'})

	return ast.WalkContinue, nil
}

var _ markdown.NodeRenderer = &ImageRenderer{}
