package container

import (
	"bytes"

	"github.com/Bornholm/amatl/pkg/markdown/directive"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// ClassPrefix is prepended to the container name to build
// the CSS class of the generated element
const ClassPrefix = "amatl-"

// semanticTags are container names rendered with the matching HTML element
// instead of a <div>
var semanticTags = map[string]struct{}{
	"article":    {},
	"aside":      {},
	"figcaption": {},
	"figure":     {},
	"footer":     {},
	"header":     {},
	"main":       {},
	"nav":        {},
	"section":    {},
}

var (
	attrNameClass = []byte("class")
	dataPrefix    = []byte("data-")
	ariaPrefix    = []byte("aria-")
)

// HTMLRenderer renders container directives as HTML elements.
//
// The element receives the "amatl-<name>" class. Global HTML attributes
// (id, class, style, role...) as well as data-* and aria-* attributes are
// rendered as is, other attributes are rendered as data-* attributes.
type HTMLRenderer struct {
}

// RegisterFuncs implements renderer.NodeRenderer.
func (r *HTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(Kind, r.render)
}

func (r *HTMLRenderer) render(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	node, ok := n.(*Node)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *container.Node, got '%T'", n)
	}

	tag := Tag(node.Name())

	if !entering {
		_, _ = w.WriteString("</" + tag + ">\n")
		return ast.WalkContinue, nil
	}

	_, _ = w.WriteString("<" + tag)

	class := ClassPrefix + node.Name()
	if value, exists := node.AttributeString(string(attrNameClass)); exists {
		if extra := directive.AttributeValueString(value); extra != "" {
			class += " " + extra
		}
	}

	writeAttribute(w, attrNameClass, class)

	for _, attr := range node.Attributes() {
		if bytes.Equal(attr.Name, attrNameClass) {
			continue
		}

		name := attr.Name
		if !html.GlobalAttributeFilter.Contains(name) && !bytes.HasPrefix(name, dataPrefix) && !bytes.HasPrefix(name, ariaPrefix) {
			name = append(append([]byte{}, dataPrefix...), name...)
		}

		writeAttribute(w, name, directive.AttributeValueString(attr.Value))
	}

	_, _ = w.WriteString(">\n")

	return ast.WalkContinue, nil
}

// Tag returns the HTML element used to render a container with the given name
func Tag(name string) string {
	if _, exists := semanticTags[name]; exists {
		return name
	}

	return "div"
}

func writeAttribute(w util.BufWriter, name []byte, value string) {
	_ = w.WriteByte(' ')
	_, _ = w.Write(name)
	_, _ = w.WriteString(`="`)
	_, _ = w.Write(util.EscapeHTML([]byte(value)))
	_ = w.WriteByte('"')
}

var _ renderer.NodeRenderer = &HTMLRenderer{}
