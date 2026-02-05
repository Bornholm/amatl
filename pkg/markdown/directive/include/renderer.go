package include

import (
	"bytes"

	"github.com/Bornholm/amatl/pkg/markdown/directive"
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
	rawURL, err := getNodeURLAttribute(node)
	if err != nil {
		panic(errors.WithStack(err))
	}

	includedSource, includedNode, exists := r.Cache.Get(rawURL)
	if !exists {
		panic(errors.Errorf("could not find source associated with path '%s'", rawURL))
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
