package attrs

import (
	"github.com/Bornholm/amatl/pkg/markdown/directive"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type NodeTransformer struct {
}

// Transform implements directive.NodeTransformer.
func (t *NodeTransformer) Transform(node *directive.Node, reader text.Reader, pc parser.Context) error {
	// Do nothing

	return nil
}

// Transform implements directive.NodeTransformer.
func (t *NodeTransformer) PostTransform(doc *ast.Document, reader text.Reader, pc parser.Context) error {
	var nextElementAttributes []ast.Attribute
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if n.Kind() != directive.KindDirective {
			if nextElementAttributes != nil {
				for _, attr := range nextElementAttributes {
					n.SetAttribute(attr.Name, attr.Value)
				}

				nextElementAttributes = nil
			}

			return ast.WalkContinue, nil
		}

		directive, ok := n.(*directive.Node)
		if !ok {
			return ast.WalkContinue, nil
		}

		if directive.DirectiveType() != Type {
			return ast.WalkContinue, nil
		}

		nextElementAttributes = directive.Attributes()

		return ast.WalkContinue, nil
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

var _ directive.NodeTransformer = &NodeTransformer{}
