package include

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/Bornholm/amatl/pkg/markdown/directive"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type NodeTransformer struct {
	BasePath string
	Cache    *SourceCache
	Parser   parser.Parser
}

// Transform implements directive.NodeTransformer.
func (t *NodeTransformer) Transform(node *directive.Node, reader text.Reader, pc parser.Context) {
	rawPath, exists := node.AttributeString("path")
	if !exists {
		return
	}

	path, ok := rawPath.(string)
	if !ok {
		return
	}

	absPath := filepath.Join(t.BasePath, string(path))

	if _, _, exists := t.Cache.Get(path); exists {
		return
	}

	includedSource, err := os.ReadFile(absPath)
	if err != nil {
		panic(errors.Wrapf(err, "could not include markdown file '%s'", absPath))
	}

	includedReader := text.NewReader(includedSource)
	includedNode := t.Parser.Parse(includedReader)

	includeDir := filepath.Dir(absPath)

	if err := t.rewriteLocalLinks(includedNode, includeDir); err != nil {
		panic(errors.Wrapf(err, "could not rewrite links of included markdown file '%s'", absPath))
	}

	t.Cache.Set(path, includedSource, includedNode)

	parent := node.Parent()
	if parent != nil && parent.Kind() == ast.KindParagraph {
		grandparent := parent.Parent()
		parent.RemoveChild(parent, node)
		grandparent.ReplaceChild(grandparent, parent, node)
	}
}

func (t *NodeTransformer) rewriteLocalLinks(root ast.Node, includeDir string) error {
	err := ast.Walk(root, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch n := node.(type) {
		case *ast.Image:
			destination := string(n.Destination)

			if isURL(destination) {
				return ast.WalkContinue, nil
			}

			if !filepath.IsAbs(destination) {
				destination = filepath.Join(includeDir, destination)

				newDestination, err := filepath.Rel(t.BasePath, destination)
				if err != nil {
					return ast.WalkStop, errors.WithStack(err)
				}

				destination = newDestination
			}

			n.Destination = []byte(destination)
		default:
			return ast.WalkContinue, nil
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

var _ directive.NodeTransformer = &NodeTransformer{}

func isURL(str string) bool {
	_, err := url.ParseRequestURI(str)
	return err == nil
}
