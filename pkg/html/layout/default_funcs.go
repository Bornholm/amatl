package layout

import (
	"bytes"
	"context"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/Masterminds/sprig/v3"
	"github.com/andybalholm/cascadia"
	"github.com/pkg/errors"
	"github.com/vincent-petithory/dataurl"
	"golang.org/x/net/html"
)

func DefaultFuncs(resolver resolver.Resolver) template.FuncMap {
	funcs := sprig.FuncMap()
	funcs["htmlQueryFirst"] = htmlQueryFirst
	funcs["htmlQueryAll"] = htmlQueryAll
	funcs["htmlSplit"] = htmlSplit
	funcs["htmlRemove"] = htmlRemove
	funcs["htmlAddAttr"] = htmlAddAttr
	funcs["htmlTextContent"] = htmlTextContent
	funcs["resolve"] = getResolveFunc(resolver)
	return funcs
}

func htmlRemove(rawHTML template.HTML, query string) (template.HTML, error) {
	buff := bytes.NewBuffer([]byte(rawHTML))

	root, err := html.Parse(buff)
	if err != nil {
		return "", errors.WithStack(err)
	}

	selector, err := cascadia.Compile(query)
	if err != nil {
		return "", errors.WithStack(err)
	}

	nodes := selector.MatchAll(root)

	for _, n := range nodes {
		n.Parent.RemoveChild(n)
	}

	rest, err := renderFragment(root)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return rest, nil
}

func htmlTextContent(rawHTML template.HTML, query string) (string, error) {
	buff := bytes.NewBuffer([]byte(rawHTML))

	root, err := html.Parse(buff)
	if err != nil {
		return "", errors.WithStack(err)
	}

	selector, err := cascadia.Compile(query)
	if err != nil {
		return "", errors.WithStack(err)
	}

	nodes := selector.MatchAll(root)

	var sb strings.Builder
	for _, n := range nodes {
		err := walk(n, func(n *html.Node) error {
			if n.Type == html.TextNode {
				sb.WriteString(n.Data)
			}
			return nil
		})
		if err != nil {
			return "", errors.WithStack(err)
		}
	}

	return sb.String(), nil
}

func htmlQueryFirst(rawHTML template.HTML, query string) (template.HTML, error) {
	buff := bytes.NewBuffer([]byte(rawHTML))

	root, err := html.Parse(buff)
	if err != nil {
		return "", errors.WithStack(err)
	}

	selector, err := cascadia.Compile(query)
	if err != nil {
		return "", errors.WithStack(err)
	}

	first := selector.MatchFirst(root)
	if first == nil {
		return "", nil
	}

	var el bytes.Buffer
	if err := html.Render(&el, first); err != nil {
		return "", errors.WithStack(err)
	}

	return template.HTML(el.String()), nil
}

func htmlQueryAll(rawHTML template.HTML, query string) ([]template.HTML, error) {
	buff := bytes.NewBuffer([]byte(rawHTML))

	root, err := html.Parse(buff)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	selector, err := cascadia.Compile(query)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	nodes := selector.MatchAll(root)

	elements := make([]template.HTML, 0, len(nodes))
	for _, n := range nodes {
		var el bytes.Buffer
		if err := html.Render(&el, n); err != nil {
			return nil, errors.WithStack(err)
		}

		elements = append(elements, template.HTML(el.String()))
	}

	return elements, nil
}

func htmlAddAttr(rawHTML template.HTML, query string, key string, value string) (template.HTML, error) {
	buff := bytes.NewBuffer([]byte(rawHTML))

	root, err := html.Parse(buff)
	if err != nil {
		return "", errors.WithStack(err)
	}

	selector, err := cascadia.Compile(query)
	if err != nil {
		return "", errors.WithStack(err)
	}

	nodes := selector.MatchAll(root)

	for _, n := range nodes {
		found := false
		for i := range n.Attr {
			if n.Attr[i].Key != key {
				continue
			}

			n.Attr[i].Val += " " + value
			found = true
			break
		}
		if !found {
			n.Attr = append(n.Attr, html.Attribute{
				Key: key,
				Val: value,
			})
		}

	}

	updated, err := renderFragment(root)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return updated, nil
}

func htmlSplit(rawHTML template.HTML, query string) ([]template.HTML, error) {
	buff := bytes.NewBuffer([]byte(rawHTML))

	root, err := html.Parse(buff)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	bodySelector, err := cascadia.Compile("body")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	elements := make([]template.HTML, 0)

	body := bodySelector.MatchFirst(root)
	if body == nil {
		return elements, nil
	}

	splitSelector, err := cascadia.Compile(query)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	appendBlock := func(block []html.Node) error {
		var el bytes.Buffer

		for _, n := range block {
			if n.Type == html.ErrorNode {
				continue
			}

			if err := html.Render(&el, &n); err != nil {
				return errors.WithStack(err)
			}
		}

		elements = append(elements, template.HTML(el.String()))

		return nil
	}

	currentBlock := make([]html.Node, 0)

	child := body.FirstChild
	for {
		if child == nil {
			break
		}

		matches := splitSelector.Match(child)

		if !matches {
			currentBlock = append(currentBlock, *child)
			child = child.NextSibling
			continue
		}

		if err := appendBlock(currentBlock); err != nil {
			return nil, errors.WithStack(err)
		}

		clear(currentBlock)

		child = child.NextSibling
	}

	if err := appendBlock(currentBlock); err != nil {
		return nil, errors.WithStack(err)
	}

	return elements, nil
}

func getResolveFunc(res resolver.Resolver) func(ctx context.Context, rawURL string, mimeTypes ...string) (template.URL, error) {
	return func(ctx context.Context, rawURL string, mimeTypes ...string) (template.URL, error) {
		path := resolver.Path(rawURL)

		reader, err := res.Resolve(ctx, path)
		if err != nil {
			return "", errors.WithStack(err)
		}

		defer reader.Close()

		data, err := io.ReadAll(reader)
		if err != nil {
			return "", errors.WithStack(err)
		}

		var mimeType string
		if len(mimeTypes) == 0 {
			mimeType = http.DetectContentType(data)
		} else {
			mimeType = mimeTypes[0]
		}

		dataURL := dataurl.New(data, mimeType)

		return template.URL(dataURL.String()), nil
	}
}

// walk calls fn on the given node and all its descendants, in document order
func walk(n *html.Node, fn func(n *html.Node) error) error {
	if err := fn(n); err != nil {
		return errors.WithStack(err)
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if err := walk(child, fn); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// renderFragment renders the content of a document parsed with html.Parse
// without the <html>, <head> and <body> elements added by the parser.
// Elements the parser moved to the <head> (i.e. a leading <style>) are kept.
func renderFragment(root *html.Node) (template.HTML, error) {
	var buf bytes.Buffer

	err := walk(root, func(n *html.Node) error {
		if n.Type != html.ElementNode || (n.Data != "head" && n.Data != "body") || n.Parent == nil || n.Parent.Data != "html" {
			return nil
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			if err := html.Render(&buf, child); err != nil {
				return errors.WithStack(err)
			}
		}

		return nil
	})
	if err != nil {
		return "", errors.WithStack(err)
	}

	return template.HTML(buf.String()), nil
}
