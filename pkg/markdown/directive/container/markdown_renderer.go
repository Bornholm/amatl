package container

import (
	"bytes"
	"strings"

	"github.com/Bornholm/amatl/pkg/markdown/directive"
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

// MarkdownRenderer renders container directives back to Markdown
type MarkdownRenderer struct {
}

// Render implements markdown.NodeRenderer.
func (*MarkdownRenderer) Render(r *markdown.Render, n ast.Node, entering bool) (ast.WalkStatus, error) {
	node, ok := n.(*Node)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *container.Node, got '%T'", n)
	}

	fence := bytes.Repeat([]byte{fenceChar}, max(node.FenceLength(), minFenceLength))

	if !entering {
		_, _ = r.Writer().Write(markdown.NewLineChar)
		_, _ = r.Writer().Write(fence)
		return ast.WalkContinue, nil
	}

	_, _ = r.Writer().Write(fence)
	_, _ = r.Writer().Write([]byte(node.Name()))

	if attrs := marshalAttributes(node.Attributes()); attrs != "" {
		_, _ = r.Writer().Write([]byte("{" + attrs + "}"))
	}

	_, _ = r.Writer().Write(markdown.NewLineChar)

	return ast.WalkContinue, nil
}

func marshalAttributes(attributes []ast.Attribute) string {
	parts := make([]string, 0, len(attributes))

	for _, attr := range attributes {
		value := directive.AttributeValueString(attr.Value)

		switch string(attr.Name) {
		case "id":
			if isShorthandSafe(value) {
				parts = append(parts, "#"+value)
				continue
			}
		case "class":
			classes := strings.Fields(value)
			if len(classes) > 0 && allShorthandSafe(classes) {
				for _, class := range classes {
					parts = append(parts, "."+class)
				}
				continue
			}
		}

		parts = append(parts, string(attr.Name)+"="+directive.QuoteAttributeValue(value))
	}

	return strings.Join(parts, " ")
}

func allShorthandSafe(values []string) bool {
	for _, v := range values {
		if !isShorthandSafe(v) {
			return false
		}
	}
	return true
}

func isShorthandSafe(value string) bool {
	if value == "" {
		return false
	}

	for _, c := range []byte(value) {
		if !(isLetter(c) || isDigit(c) || c == '-' || c == '_' || c == ':') {
			return false
		}
	}

	return true
}

var _ markdown.NodeRenderer = &MarkdownRenderer{}
