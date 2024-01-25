package markdown

import (
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type NodeRendererHookFunc func(r *Render, node ast.Node, entering bool)

type hookedNodeRenderer struct {
	before   NodeRendererHookFunc
	renderer NodeRenderer
	after    NodeRendererHookFunc
}

// Render implements NodeRenderer.
func (hr *hookedNodeRenderer) Render(r *Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if hr.before != nil {
		hr.before(r, node, entering)
	}

	status, err := hr.renderer.Render(r, node, entering)
	if err != nil {
		return status, errors.WithStack(err)
	}

	if hr.after != nil {
		hr.after(r, node, entering)
	}

	return status, err
}

var _ NodeRenderer = &hookedNodeRenderer{}

func WithBeforeRender(renderer NodeRenderer, before NodeRendererHookFunc) NodeRenderer {
	return &hookedNodeRenderer{
		before:   before,
		renderer: renderer,
	}
}

func WithAfterRender(renderer NodeRenderer, after NodeRendererHookFunc) NodeRenderer {
	return &hookedNodeRenderer{
		after:    after,
		renderer: renderer,
	}
}
