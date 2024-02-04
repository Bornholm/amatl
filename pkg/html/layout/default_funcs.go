package layout

import (
	"bytes"
	"html/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/andybalholm/cascadia"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

func DefaultFuncs() template.FuncMap {
	funcs := sprig.FuncMap()
	funcs["htmlQueryAll"] = htmlQueryAll
	funcs["htmlSplit"] = htmlSplit
	return funcs
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
