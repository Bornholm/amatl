package include

import (
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

	t.Cache.Set(path, includedSource, includedNode)

	parent := node.Parent()
	if parent != nil && parent.Kind() == ast.KindParagraph {
		grandparent := parent.Parent()
		parent.RemoveChild(parent, node)
		grandparent.ReplaceChild(grandparent, parent, node)
	}
}

var _ directive.NodeTransformer = &NodeTransformer{}
