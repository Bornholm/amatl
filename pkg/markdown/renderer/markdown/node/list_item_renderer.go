package node

import (
	"bytes"

	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type ListItemRenderer struct {
}

// Render implements NodeRenderer.
func (*ListItemRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	listItem, ok := node.(*ast.ListItem)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.String, got '%T'", node)
	}

	if entering {
		liMarker := listItemMarkerChars(listItem)
		_, _ = r.Writer().Write(liMarker)
		if r.Renderer().ListIndentStyle() == markdown.ListIndentUniform &&
			// We can use 4 spaces for indentation only if
			// that would still qualify as part of the list
			// item text. e.g., given "123. foo",
			// for content to be part of that list item,
			// it must be indented 5 spaces.
			//
			//	123. foo
			//
			//	     bar
			len(liMarker) <= len(markdown.FourSpacesChars) {
			r.Writer().PushIndent(markdown.FourSpacesChars)
		} else {
			r.Writer().PushIndent(bytes.Repeat(markdown.SpaceChar, len(liMarker)))
		}
	} else {
		if listItem.NextSibling() != nil && listItem.NextSibling().Kind() == ast.KindListItem {
			// Newline after list item.
			_, _ = r.Writer().Write(markdown.NewLineChar)
		}
		r.Writer().PopIndent()
	}

	return ast.WalkContinue, nil
}

var _ markdown.NodeRenderer = &ListItemRenderer{}
