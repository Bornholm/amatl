package include

import (
	"bytes"

	"forge.cadoles.com/wpetit/amatl/pkg/markdown/directive"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type NodeRenderer struct {
	Cache    *SourceCache
	Renderer renderer.Renderer
}

// Render implements directive.NodeRenderer.
func (r *NodeRenderer) Render(writer util.BufWriter, source []byte, node *directive.Node) {
	rawPath, exists := node.AttributeString("path")
	if !exists {
		panic(errors.Errorf("could not find 'path' attribute in directive '%s'", node.Kind()))
	}

	path, ok := rawPath.(string)
	if !ok {
		panic(errors.Errorf("unexpected type '%T' for 'path' attribute", rawPath))
	}

	includedSource, includedNode, exists := r.Cache.Get(path)
	if !exists {
		panic(errors.Errorf("could not find source associated with path '%s'", path))
	}

	var buff bytes.Buffer

	if err := r.Renderer.Render(&buff, includedSource, includedNode); err != nil {
		panic(errors.Wrap(err, "could not convert included source"))
	}

	if _, err := writer.Write(buff.Bytes()); err != nil {
		panic(errors.WithStack(err))
	}
}

var _ directive.NodeRenderer = &NodeRenderer{}
