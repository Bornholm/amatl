package directive

import (
	"fmt"
	"strings"

	"forge.cadoles.com/wpetit/amatl/pkg/markdown/renderer/markdown"
	"github.com/yuin/goldmark/ast"
)

type MarkdownDirectiveRenderer interface {
	Render(r *markdown.Render, directive *Node, entering bool) (ast.WalkStatus, error)
}

type markdownDirectiveRenderer struct {
	render MarkdownDirectiveRendererFunc
}

func (t *markdownDirectiveRenderer) Render(r *markdown.Render, directive *Node, entering bool) (ast.WalkStatus, error) {
	return t.render(r, directive, entering)
}

type MarkdownDirectiveRendererFunc func(r *markdown.Render, directive *Node, entering bool) (ast.WalkStatus, error)

type MarkdownNodeRenderer struct {
	renderers map[Type]MarkdownDirectiveRenderer
}

// Render implements markdown.NodeRenderer.
func (mr *MarkdownNodeRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	directive, ok := node.(*Node)
	if !ok {
		return ast.WalkStop, fmt.Errorf("unexpected node %T, expected *directive.Node", node)
	}

	renderer, exists := mr.renderers[directive.DirectiveType()]
	if !exists {
		return mr.renderDefault(r, directive, entering)
	}

	return renderer.Render(r, directive, entering)
}

func (mr *MarkdownNodeRenderer) renderDefault(r *markdown.Render, directive *Node, entering bool) (ast.WalkStatus, error) {
	str := fmt.Sprintf(":%s{%s}", directive.DirectiveType(), marshalAttributes(directive.Attributes()))

	_, _ = r.Writer().Write(markdown.NewLineChar)
	_, _ = r.Writer().Write(markdown.NewLineChar)
	_, _ = r.Writer().Write([]byte(str))

	return ast.WalkContinue, nil
}

func marshalAttributes(attributes []ast.Attribute) string {
	var sb strings.Builder

	for idx, attr := range attributes {
		if idx != 0 {
			sb.WriteString(" ")
		}
		sb.Write(attr.Name)
		sb.WriteString("=")
		sb.WriteRune('"')
		sb.WriteString(fmt.Sprintf("%v", attr.Value))
		sb.WriteRune('"')
	}

	return sb.String()
}

func NewMarkdownNodeRenderer(funcs ...MarkdownNodeRendererOptionFunc) *MarkdownNodeRenderer {
	opts := NewMarkdownNodeRendererOptions(funcs...)
	return &MarkdownNodeRenderer{
		renderers: opts.Renderers,
	}
}

type MarkdownNodeRendererOptions struct {
	Renderers map[Type]MarkdownDirectiveRenderer
}

type MarkdownNodeRendererOptionFunc func(opts *MarkdownNodeRendererOptions)

func NewMarkdownNodeRendererOptions(funcs ...MarkdownNodeRendererOptionFunc) *MarkdownNodeRendererOptions {
	opts := &MarkdownNodeRendererOptions{
		Renderers: make(map[Type]MarkdownDirectiveRenderer),
	}

	for _, fn := range funcs {
		fn(opts)
	}

	return opts
}

func WithMarkdownDirectiveRenderer(directiveType Type, renderer MarkdownDirectiveRenderer) MarkdownNodeRendererOptionFunc {
	return func(opts *MarkdownNodeRendererOptions) {
		opts.Renderers[directiveType] = renderer
	}
}

func WithMarkdownDirectiveRendererFunc(directiveType Type, render MarkdownDirectiveRendererFunc) MarkdownNodeRendererOptionFunc {
	return func(opts *MarkdownNodeRendererOptions) {
		opts.Renderers[directiveType] = &markdownDirectiveRenderer{render}
	}
}

var _ markdown.NodeRenderer = &MarkdownNodeRenderer{}
