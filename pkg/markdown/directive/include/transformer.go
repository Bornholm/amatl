package include

import (
	"io"
	"net/url"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/Bornholm/amatl/pkg/markdown/directive"
	"github.com/Bornholm/amatl/pkg/markdown/directive/container"
	"github.com/Bornholm/amatl/pkg/markdown/selector"
	"github.com/Bornholm/amatl/pkg/pipeline"
	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/Bornholm/amatl/pkg/transform"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type NodeTransformer struct {
	Parser     parser.Parser
	SourcePath resolver.Path
}

// Transform implements directive.NodeTransformer.
func (t *NodeTransformer) Transform(node *directive.Node, reader text.Reader, pc parser.Context) error {
	sourcePath := getSourcePath(pc, t.SourcePath)

	resourcePath, _, err := parseNodeURLAttribute(sourcePath, node)
	if err != nil {
		return errors.Wrapf(err, "could not parse required attribute on directive '%s'", node.DirectiveType())
	}

	chain := getIncludeChain(pc, sourcePath)
	if slices.Contains(chain, includeChainKey(resourcePath)) {
		return errors.Errorf("include cycle: %s -> %s", strings.Join(chain, " -> "), includeChainKey(resourcePath))
	}

	shiftHeadings, err := getNodeShiftHeadingsAttribute(node)
	if err != nil {
		shiftHeadings = 0
	}

	fromHeadings, err := getNodeFromHeadingsAttribute(node)
	if err != nil {
		fromHeadings = 0
	}

	ctx, err := pipeline.FromParserContext(pc)
	if err != nil {
		return errors.WithStack(err)
	}

	resourceReader, err := resolver.Resolve(ctx, resourcePath.String())
	if err != nil {
		return errors.Wrapf(err, "could not resolve resource '%s'", resourcePath)
	}

	defer func() {
		if err := resourceReader.Close(); err != nil {
			panic(errors.Wrapf(err, "could not close resource '%s'", resourcePath))
		}
	}()

	transformed := transform.NewNewlineReader(resourceReader)

	includedSource, err := io.ReadAll(transformed)
	if err != nil {
		return errors.Wrapf(err, "could not read markdown resource '%s'", resourcePath)
	}

	setIncludedSource(node, includedSource)

	includedReader := text.NewReader(includedSource)

	sourceDir, err := resourcePath.Dir().Abs()
	if err != nil {
		return errors.Wrapf(err, "could not retrieve absolute path to resource '%s'", resourcePath)
	}

	includeCtx := resolver.WithWorkDir(ctx, sourceDir)
	includePC := pipeline.WithContext(includeCtx, parser.NewContext())

	setSourcePath(includePC, resourcePath)
	setIncludeChain(includePC, append(slices.Clone(chain), includeChainKey(resourcePath)))

	includedNode := t.Parser.Parse(includedReader, parser.WithContext(includePC))

	selectAttr, _ := getNodeSelectAttribute(node)
	if selectAttr != "" {
		sel, err := selector.Parse(selectAttr)
		if err != nil {
			return errors.Wrapf(err, "could not parse 'select' attribute on :include directive for '%s'", resourcePath)
		}
		matchedNodes := sel.MatchTopLevel(includedNode, includedSource)
		includedNode = buildFilteredDocument(matchedNodes)
	}

	if err := t.excludeSections(includedNode, fromHeadings); err != nil {
		return errors.Wrapf(err, "could not exclude sections of included markdown resource '%s'", resourcePath)
	}

	if err := t.rewriteRelativeLinks(includedNode, resourcePath); err != nil {
		return errors.Wrapf(err, "could not rewrite links of included markdown resource '%s'", resourcePath)
	}

	if err := t.shiftHeadings(includedNode, shiftHeadings); err != nil {
		return errors.Wrapf(err, "could not shift headings of included markdown resource '%s'", resourcePath)
	}

	setIncludedNode(node, includedNode)

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
				return ast.WalkSkipChildren, nil
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

func (t *NodeTransformer) rewriteRelativeLinks(root ast.Node, basePath resolver.Path) error {
	dirPath := basePath.Dir()

	err := ast.Walk(root, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch n := node.(type) {
		case *ast.Link:
			destination := string(n.Destination)

			if isURL(destination) {
				return ast.WalkContinue, nil
			}

			if !filepath.IsAbs(destination) {
				destination = dirPath.JoinPath(destination).String()
			}

			n.Destination = []byte(destination)
		case *ast.Image:
			destination := string(n.Destination)

			if isURL(destination) {
				return ast.WalkContinue, nil
			}

			if !filepath.IsAbs(destination) {
				destination = dirPath.JoinPath(destination).String()
			}

			n.Destination = []byte(destination)
		case *container.Node:
			value, exists := n.AttributeString(container.AttrNameBackground)
			if !exists {
				return ast.WalkContinue, nil
			}

			destination := directive.AttributeValueString(value)

			if destination == "" || isURL(destination) || filepath.IsAbs(destination) {
				return ast.WalkContinue, nil
			}

			n.SetAttributeString(container.AttrNameBackground, []byte(dirPath.JoinPath(destination).String()))
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
	return err == nil || strings.HasPrefix(str, "#")
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

func parseNodeURLAttribute(basePath resolver.Path, node ast.Node) (resolver.Path, string, error) {
	rawURL, err := getNodeURLAttribute(node)
	if err != nil {
		return "", "", errors.WithStack(err)
	}

	baseDir := basePath.Dir()

	return baseDir.Join(resolver.Path(rawURL)), rawURL, nil
}

var contextKeySourcePath = parser.NewContextKey()

var contextKeyIncludeChain = parser.NewContextKey()

// getIncludeChain returns the documents being included, from the root
// document to the one holding the directive
func getIncludeChain(ctx parser.Context, sourcePath resolver.Path) []string {
	chain, ok := ctx.Get(contextKeyIncludeChain).([]string)
	if !ok || len(chain) == 0 {
		return []string{includeChainKey(sourcePath)}
	}

	return chain
}

func setIncludeChain(ctx parser.Context, chain []string) {
	ctx.Set(contextKeyIncludeChain, chain)
}

// includeChainKey returns the absolute form of the path, so that the same
// document reached through different relative paths is recognized
func includeChainKey(path resolver.Path) string {
	if abs, err := path.Abs(); err == nil {
		return abs.String()
	}

	return path.String()
}

func getSourcePath(ctx parser.Context, defaultSourcePath resolver.Path) resolver.Path {
	sourcePath, ok := ctx.Get(contextKeySourcePath).(resolver.Path)
	if !ok || sourcePath == "" {
		return defaultSourcePath
	}

	return sourcePath
}

func setSourcePath(ctx parser.Context, path resolver.Path) {
	ctx.Set(contextKeySourcePath, path)
}

const attrNameSelect = "select"

func getNodeSelectAttribute(node ast.Node) (string, bool) {
	selectAttrValue, exists := node.AttributeString(attrNameSelect)
	if !exists {
		return "", false
	}
	selectStr, ok := selectAttrValue.(string)
	if !ok {
		return "", false
	}
	return selectStr, true
}

// buildFilteredDocument creates a new ast.Document containing only the given nodes.
func buildFilteredDocument(nodes []ast.Node) *ast.Document {
	doc := ast.NewDocument()
	for _, node := range nodes {
		parent := node.Parent()
		if parent != nil {
			parent.RemoveChild(parent, node)
		}
		doc.AppendChild(doc, node)
	}
	return doc
}
