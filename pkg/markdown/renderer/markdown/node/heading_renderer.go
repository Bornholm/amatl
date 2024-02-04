package node

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/mattn/go-runewidth"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type HeadingRenderer struct {
}

// Render implements NodeRenderer.
func (hr *HeadingRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	heading, ok := node.(*ast.Heading)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *ast.Heading, got '%T'", node)
	}

	if !entering {
		return ast.WalkContinue, nil
	}

	// Render it straight away. No nested headings are supported and we expect
	// headings to have limited content, so limit WALK.
	if err := hr.renderHeading(r, heading); err != nil {
		return ast.WalkStop, fmt.Errorf("rendering heading: %w", err)
	}
	return ast.WalkSkipChildren, nil
}

func (hr *HeadingRenderer) renderHeading(r *markdown.Render, node *ast.Heading) error {
	underlineHeading := false
	if r.Renderer().UnderlineHeadings() {
		underlineHeading = node.Level <= 2
	}

	if !underlineHeading {
		r.Writer().Write(bytes.Repeat([]byte{'#'}, node.Level))
		r.Writer().Write(markdown.SpaceChar)
	}

	var headBuf bytes.Buffer
	headBuf.Reset()

	for n := node.FirstChild(); n != nil; n = n.NextSibling() {
		if err := ast.Walk(n, func(inner ast.Node, entering bool) (ast.WalkStatus, error) {
			if entering {
				if err := ast.Walk(inner, r.Renderer().NewRender(&headBuf, r.Source()).RenderNode); err != nil {
					return ast.WalkStop, err
				}
			}
			return ast.WalkSkipChildren, nil
		}); err != nil {
			return err
		}
	}
	a := node.Attributes()
	sort.SliceStable(a, func(i, j int) bool {
		switch {
		case bytes.Equal(a[i].Name, []byte("id")):
			return true
		case bytes.Equal(a[j].Name, []byte("id")):
			return false
		case bytes.Equal(a[i].Name, []byte("class")):
			return true
		case bytes.Equal(a[j].Name, []byte("class")):
			return false
		}
		return bytes.Compare(a[i].Name, a[j].Name) == -1
	})

	// hAttr := []string{}
	// for _, attr := range node.Attributes() {
	// 	switch string(attr.Name) {
	// 	case "id":
	// 		hAttr = append(hAttr, fmt.Sprintf("#%s", attr.Value))
	// 	case "class":
	// 		hAttr = append(hAttr, strings.ReplaceAll(fmt.Sprintf(".%s", attr.Value), " ", " ."))
	// 	default:
	// 		if attr.Value == nil {
	// 			hAttr = append(hAttr, string(attr.Name))
	// 			continue
	// 		}
	// 		hAttr = append(hAttr, fmt.Sprintf("%s=%s", string(attr.Name), attr.Value))
	// 	}
	// }
	// if len(hAttr) != 0 {
	// 	_, _ = fmt.Fprintf(&headBuf, " {%s}", strings.Join(hAttr, " "))
	// }

	_, _ = r.Writer().Write(headBuf.Bytes())

	if underlineHeading {
		width := runewidth.StringWidth(headBuf.String())

		_, _ = r.Writer().Write(markdown.NewLineChar)

		switch node.Level {
		case 1:
			r.Writer().Write(bytes.Repeat(markdown.Heading1UnderlineChar, width))
		case 2:
			r.Writer().Write(bytes.Repeat(markdown.Heading2UnderlineChar, width))
		}
	}

	return nil
}

var _ markdown.NodeRenderer = &HeadingRenderer{}
