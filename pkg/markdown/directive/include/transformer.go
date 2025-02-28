package include

import (
	"context"
	"io"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/Bornholm/amatl/pkg/markdown/directive"
	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type NodeTransformer struct {
	Cache     *SourceCache
	Parser    parser.Parser
	SourceURL *url.URL
}

// Transform implements directive.NodeTransformer.
func (t *NodeTransformer) Transform(node *directive.Node, reader text.Reader, pc parser.Context) error {
	sourceURL := getSourceURL(pc, t.SourceURL)

	rawURL, resourceURL, err := parseNodeURLAttribute(sourceURL, node)
	if err != nil {
		return errors.Wrapf(err, "could not parse required attribute on directive '%s'", node.DirectiveType())
	}

	shiftHeadings, err := getNodeShiftHeadingsAttribute(node)
	if err != nil {
		shiftHeadings = 0
	}

	fromHeadings, err := getNodeFromHeadingsAttribute(node)
	if err != nil {
		fromHeadings = 0
	}

	if _, _, exists := t.Cache.Get(resourceURL.String()); exists {
		return nil
	}

	resourceReader, err := resolver.Resolve(context.Background(), resourceURL)
	if err != nil {
		return errors.Wrapf(err, "could not resolve resource '%s'", resourceURL)
	}

	defer func() {
		if err := resourceReader.Close(); err != nil {
			panic(errors.Wrapf(err, "could not close resource '%s'", resourceURL))
		}
	}()

	includedSource, err := io.ReadAll(resourceReader)
	if err != nil {
		return errors.Wrapf(err, "could not read markdown resource '%s'", resourceURL)
	}

	setIncludedSource(node, includedSource)

	includedReader := text.NewReader(includedSource)

	ctx := parser.NewContext()
	setSourceURL(ctx, resourceURL)

	includedNode := t.Parser.Parse(includedReader, parser.WithContext(ctx))

	if err := t.excludeSections(includedNode, fromHeadings); err != nil {
		return errors.Wrapf(err, "could not exclude sections of included markdown resource '%s'", resourceURL)
	}

	if err := t.rewriteRelativeLinks(includedNode, resourceURL); err != nil {
		return errors.Wrapf(err, "could not rewrite links of included markdown resource '%s'", resourceURL)
	}

	if err := t.shiftHeadings(includedNode, shiftHeadings); err != nil {
		return errors.Wrapf(err, "could not shift headings of included markdown resource '%s'", resourceURL)
	}

	setIncludedNode(node, includedNode)
	t.Cache.Set(rawURL, includedSource, includedNode)

	parent := node.Parent()
	if parent != nil && parent.Kind() == ast.KindParagraph {
		grandparent := parent.Parent()
		parent.RemoveChild(parent, node)
		grandparent.ReplaceChild(grandparent, parent, node)
	}

	return nil
}

func (t *NodeTransformer) excludeSections(root ast.Node, minLevel int) error {
	currentLevel := 0
	err := ast.Walk(root, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if heading, ok := node.(*ast.Heading); ok {
			currentLevel = heading.Level
		}

		if currentLevel < minLevel {
			parent := node.Parent()
			if parent != nil {
				parent.RemoveChild(parent, node)
			}
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (t *NodeTransformer) shiftHeadings(root ast.Node, shift int) error {
	err := ast.Walk(root, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch n := node.(type) {
		case *ast.Heading:
			newLevel := n.Level + shift
			if newLevel > 6 {
				newLevel = 6
			}
			n.Level = newLevel
		default:
			return ast.WalkContinue, nil
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (t *NodeTransformer) rewriteRelativeLinks(root ast.Node, baseURL *url.URL) error {
	dirUrl := urlDir(baseURL)
	err := ast.Walk(root, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch n := node.(type) {
		case *ast.Image:
			destination := string(n.Destination)

			if isURL(destination) {
				return ast.WalkContinue, nil
			}

			if !filepath.IsAbs(destination) {
				newDestination := dirUrl.JoinPath(destination)
				destination = newDestination.String()
			}

			n.Destination = []byte(destination)
		default:
			return ast.WalkContinue, nil
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

var _ directive.NodeTransformer = &NodeTransformer{}

func isURL(str string) bool {
	_, err := url.ParseRequestURI(str)
	return err == nil
}

const attrNameUrl = "url"

func getNodeURLAttribute(node ast.Node) (string, error) {
	urlAttrValue, exists := node.AttributeString(attrNameUrl)
	if !exists {
		return "", errors.Errorf("attribute '%s' not found", attrNameUrl)
	}

	rawURL, ok := urlAttrValue.(string)
	if !ok {
		return "", errors.Errorf("unexpected value type '%T' for '%s' attribute", urlAttrValue, attrNameUrl)
	}

	return rawURL, nil
}

const attrNameShiftHeadings = "shiftHeadings"

func getNodeShiftHeadingsAttribute(node ast.Node) (int, error) {
	shiftHeadingsAttrValue, exists := node.AttributeString(attrNameShiftHeadings)
	if !exists {
		return 0, errors.Errorf("attribute '%s' not found", attrNameShiftHeadings)
	}

	rawShiftHeadings, ok := shiftHeadingsAttrValue.(string)
	if !ok {
		return 0, errors.Errorf("unexpected value type '%T' for '%s' attribute", shiftHeadingsAttrValue, attrNameShiftHeadings)
	}

	shiftHeadings, err := strconv.ParseInt(rawShiftHeadings, 10, 64)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return int(shiftHeadings), nil
}

const attrNameFromHeadings = "fromHeadings"

func getNodeFromHeadingsAttribute(node ast.Node) (int, error) {
	fromHeadingsAttrValue, exists := node.AttributeString(attrNameFromHeadings)
	if !exists {
		return 0, errors.Errorf("attribute '%s' not found", attrNameFromHeadings)
	}

	rawFromHeadings, ok := fromHeadingsAttrValue.(string)
	if !ok {
		return 0, errors.Errorf("unexpected value type '%T' for '%s' attribute", attrNameFromHeadings, attrNameShiftHeadings)
	}

	fromHeadings, err := strconv.ParseInt(rawFromHeadings, 10, 64)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return int(fromHeadings), nil
}

func parseNodeURLAttribute(baseURL *url.URL, node ast.Node) (string, *url.URL, error) {
	rawURL, err := getNodeURLAttribute(node)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	baseDir := urlDir(baseURL)

	var resourceURL *url.URL

	switch {
	case !isURL(rawURL) && !filepath.IsAbs(rawURL):
		resourceURL = baseDir.JoinPath(rawURL)
	case isURL(rawURL):
		resourceURL, err = url.Parse(rawURL)
		if err != nil {
			return "", nil, errors.Wrapf(err, "could not parse resource url '%s'", rawURL)
		}

	default:
		resourceURL, err = url.Parse(rawURL)
		if err != nil {
			return "", nil, errors.Wrapf(err, "could not parse resource url '%s'", rawURL)
		}

		resourceURL = baseDir.JoinPath(resourceURL.Path)
	}

	return rawURL, resourceURL, nil
}

var contextKeySourceURL = parser.NewContextKey()

func getSourceURL(ctx parser.Context, defaultSourceURL *url.URL) *url.URL {
	sourceURL, ok := ctx.Get(contextKeySourceURL).(*url.URL)
	if !ok || sourceURL == nil {
		return defaultSourceURL
	}

	return sourceURL
}

func setSourceURL(ctx parser.Context, url *url.URL) {
	ctx.Set(contextKeySourceURL, url)
}

func urlDir(src *url.URL) *url.URL {
	dir, _ := url.Parse(src.String())
	dir.Path = filepath.Dir(src.Path)
	return dir
}
