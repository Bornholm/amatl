package markdown

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/mattn/go-runewidth"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

type HeadingRenderer struct {
}

// Render implements NodeRenderer.
func (hr *HeadingRenderer) Render(r *Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
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

func (hr *HeadingRenderer) renderHeading(r *Render, node *ast.Heading) error {
	underlineHeading := false
	if r.mr.underlineHeadings {
		underlineHeading = node.Level <= 2
	}

	if !underlineHeading {
		r.w.Write(bytes.Repeat([]byte{'#'}, node.Level))
		r.w.Write(SpaceChar)
	}

	var headBuf bytes.Buffer
	headBuf.Reset()

	for n := node.FirstChild(); n != nil; n = n.NextSibling() {
		if err := ast.Walk(n, func(inner ast.Node, entering bool) (ast.WalkStatus, error) {
			if entering {
				if err := ast.Walk(inner, r.mr.newRender(&headBuf, r.source).renderNode); err != nil {
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

	_, _ = r.w.Write(headBuf.Bytes())

	if underlineHeading {
		width := runewidth.StringWidth(headBuf.String())

		_, _ = r.w.Write(NewLineChar)

		switch node.Level {
		case 1:
			r.w.Write(bytes.Repeat(Heading1UnderlineChar, width))
		case 2:
			r.w.Write(bytes.Repeat(Heading2UnderlineChar, width))
		}
	}

	return nil
}

var _ NodeRenderer = &HeadingRenderer{}
