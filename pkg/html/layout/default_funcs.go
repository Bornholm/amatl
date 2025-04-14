package layout

import (
	"bytes"
	"context"
	"html/template"
	"io"
	"net/http"
	"net/url"
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
		var el bytes.Buffer
		if err := html.Render(&el, n); err != nil {
			return "", errors.WithStack(err)
		}

		n.Parent.RemoveChild(n)
	}

	var rest bytes.Buffer
	if err := html.Render(&rest, root); err != nil {
		return "", errors.WithStack(err)
	}

	return template.HTML(rest.String()), nil
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
		for _, attr := range n.Attr {
			if attr.Key != key {
				continue
			}

			attr.Val += " " + value
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

	var updated bytes.Buffer
	if err := html.Render(&updated, root); err != nil {
		return "", errors.WithStack(err)
	}

	return template.HTML(updated.String()), nil
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

func getResolveFunc(resolver resolver.Resolver) func(ctx context.Context, rawURL string, mimeTypes ...string) (template.URL, error) {
	return func(ctx context.Context, rawURL string, mimeTypes ...string) (template.URL, error) {
		url, err := url.Parse(rawURL)
		if err != nil {
			return "", errors.WithStack(err)
		}

		reader, err := resolver.Resolve(ctx, url)
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

func walk(n *html.Node, fn func(n *html.Node) error) error {
	if n.FirstChild != nil {
		if err := walk(n.FirstChild, fn); err != nil {
			return errors.WithStack(err)
		}
	}

	if n.NextSibling != nil {
		if err := walk(n.NextSibling, fn); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
