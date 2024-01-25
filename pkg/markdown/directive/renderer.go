package directive

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type NodeRenderer interface {
	Render(writer util.BufWriter, source []byte, node *Node)
}

type nodeRenderer struct {
	render NodeRendererFunc
}

func (t *nodeRenderer) Render(writer util.BufWriter, source []byte, node *Node) {
	t.render(writer, source, node)
}

type NodeRendererFunc func(writer util.BufWriter, source []byte, node *Node)

type Renderer struct {
	renderers map[Type]NodeRenderer
}

// Render implements renderer.Renderer.
func (r *Renderer) Render(writer util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	directive, ok := node.(*Node)
	if !ok {
		return ast.WalkStop, fmt.Errorf("unexpected node %T, expected *directive.Node", node)
	}

	renderer, exists := r.renderers[directive.DirectiveType()]
	if !exists {
		return ast.WalkSkipChildren, nil
	}

	renderer.Render(writer, source, directive)

	return ast.WalkContinue, nil
}

func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindDirective, r.Render)
}

func NewRenderer(funcs ...RendererOptionFunc) *Renderer {
	opts := NewRendererOptions(funcs...)
	return &Renderer{
		renderers: opts.Renderers,
	}
}

type RendererOptions struct {
	Renderers map[Type]NodeRenderer
}

type RendererOptionFunc func(opts *RendererOptions)

func NewRendererOptions(funcs ...RendererOptionFunc) *RendererOptions {
	opts := &RendererOptions{
		Renderers: make(map[Type]NodeRenderer),
	}

	for _, fn := range funcs {
		fn(opts)
	}

	return opts
}

func WithRenderer(directiveType Type, Renderer NodeRenderer) RendererOptionFunc {
	return func(opts *RendererOptions) {
		opts.Renderers[directiveType] = Renderer
	}
}

func WithRendererFunc(directiveType Type, transform NodeRendererFunc) RendererOptionFunc {
	return func(opts *RendererOptions) {
		opts.Renderers[directiveType] = &nodeRenderer{transform}
	}
}

var _ renderer.NodeRenderer = &Renderer{}
