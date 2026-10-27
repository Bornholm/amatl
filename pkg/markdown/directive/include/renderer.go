package include

import (
	"bytes"

	"github.com/Bornholm/amatl/pkg/markdown/directive"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type NodeRenderer struct {
	Renderer renderer.Renderer
}

// Render implements directive.NodeRenderer.
func (r *NodeRenderer) Render(writer util.BufWriter, source []byte, node *directive.Node) {
	includedSource, includedNode, err := includedContent(node)
	if err != nil {
		panic(errors.WithStack(err))
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
