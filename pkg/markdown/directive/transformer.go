package directive

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type NodeTransformer interface {
	Transform(node *Node, reader text.Reader, pc parser.Context)
}

type nodeTransformer struct {
	transform NodeTransformerFunc
}

func (t *nodeTransformer) Transform(node *Node, reader text.Reader, pc parser.Context) {
	t.transform(node, reader, pc)
}

type NodeTransformerFunc func(node *Node, reader text.Reader, pc parser.Context)

type Transformer struct {
	transformers map[Type]NodeTransformer
}

// Transform implements parser.ASTTransformer.
func (t *Transformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		directive, ok := n.(*Node)
		if !ok {
			return ast.WalkContinue, nil
		}

		directiveType := directive.DirectiveType()
		transformer, exists := t.transformers[directiveType]
		if !exists {
			return ast.WalkContinue, nil
		}

		transformer.Transform(directive, reader, pc)

		return ast.WalkSkipChildren, nil
	})
}

func NewTransformer(funcs ...TransformerOptionFunc) *Transformer {
	opts := NewTransformerOptions(funcs...)
	return &Transformer{
		transformers: opts.Transformers,
	}
}

type TransformerOptions struct {
	Transformers map[Type]NodeTransformer
}

type TransformerOptionFunc func(opts *TransformerOptions)

func NewTransformerOptions(funcs ...TransformerOptionFunc) *TransformerOptions {
	opts := &TransformerOptions{
		Transformers: make(map[Type]NodeTransformer),
	}

	for _, fn := range funcs {
		fn(opts)
	}

	return opts
}

func WithTransformer(directiveType Type, transformer NodeTransformer) TransformerOptionFunc {
	return func(opts *TransformerOptions) {
		opts.Transformers[directiveType] = transformer
	}
}

func WithTransformerFunc(directiveType Type, transform NodeTransformerFunc) TransformerOptionFunc {
	return func(opts *TransformerOptions) {
		opts.Transformers[directiveType] = &nodeTransformer{transform}
	}
}

var _ parser.ASTTransformer = &Transformer{}
