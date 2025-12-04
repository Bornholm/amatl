package toc

import (
	"math"
	"strconv"

	"github.com/Bornholm/amatl/pkg/markdown/directive"
	"github.com/Bornholm/amatl/pkg/markdown/directive/include"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type tocItem struct {
	Level    int
	Label    []byte
	ID       []byte
	Children []*tocItem
	Parent   *tocItem
}

func (i *tocItem) FirstAncestor(level int) *tocItem {
	if i.Level < level {
		return i
	} else if i.Parent != nil {
		return i.Parent.FirstAncestor(level)
	}
	return nil
}

type NodeTransformer struct {
}

// Transform implements directive.NodeTransformer.
func (t *NodeTransformer) Transform(node *directive.Node, reader text.Reader, pc parser.Context) error {
	// Do nothing

	return nil
}

// Transform implements directive.NodeTransformer.
func (t *NodeTransformer) PostTransform(doc *ast.Document, reader text.Reader, pc parser.Context) error {
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if n.Kind() != directive.KindDirective {
			return ast.WalkContinue, nil
		}

		directive, ok := n.(*directive.Node)
		if !ok {
			return ast.WalkContinue, nil
		}

		if directive.DirectiveType() != Type {
			return ast.WalkContinue, nil
		}

		minLevel, err := getNodeMinLevelAttribute(directive)
		if err != nil {
			return ast.WalkStop, errors.Wrapf(err, "could not parse '%s' attribute on directive '%s'", attrNameMinLevel, directive.DirectiveType())
		}

		maxLevel, err := getNodeMaxLevelAttribute(directive)
		if err != nil {
			return ast.WalkStop, errors.Wrapf(err, "could not parse '%s' attribute on directive '%s'", attrNameMaxLevel, directive.DirectiveType())
		}

		tree, err := buildTree(doc, reader, minLevel, maxLevel)
		if err != nil {
			return ast.WalkStop, errors.WithStack(err)
		}

		toc, err := treeToNode(tree)
		if err != nil {
			return ast.WalkStop, errors.WithStack(err)
		}

		toc.SetAttribute([]byte("class"), "amatl-toc")

		parent := n.Parent()
		parent.ReplaceChild(parent, n, toc)

		return ast.WalkContinue, nil
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

var _ directive.NodeTransformer = &NodeTransformer{}

const attrNameMinLevel = "minLevel"

func treeToNode(tree []*tocItem) (ast.Node, error) {
	list := ast.NewList('-')

	for _, item := range tree {
		listItem := ast.NewListItem(0)

		if t := item.Label; len(t) > 0 {
			title := ast.NewString(t)
			title.SetRaw(true)
			if len(item.ID) > 0 {
				link := ast.NewLink()
				link.Destination = append([]byte("#"), item.ID...)
				link.AppendChild(link, title)
				listItem.AppendChild(listItem, link)
			} else {
				listItem.AppendChild(listItem, title)
			}
		}

		children, err := treeToNode(item.Children)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		listItem.AppendChild(listItem, children)
		list.AppendChild(list, listItem)
	}

	return list, nil
}

func buildTree(root ast.Node, reader text.Reader, minLevel int, maxLevel int) ([]*tocItem, error) {
	stack := []*tocItem{}

	var lastInserted *tocItem

	appendHeading := func(heading *ast.Heading) {
		label := heading.Text(reader.Source())

		if heading.Level < minLevel || heading.Level > maxLevel {
			return
		}

		newItem := &tocItem{
			Level:    heading.Level,
			Label:    label,
			Children: []*tocItem{},
		}

		if id, ok := heading.AttributeString("id"); ok {
			newItem.ID, _ = id.([]byte)
		}

		if lastInserted != nil {
			if lastInserted.Level < newItem.Level {
				lastInserted.Children = append(lastInserted.Children, newItem)
				newItem.Parent = lastInserted
			} else {
				firstAncestor := lastInserted.FirstAncestor(newItem.Level)
				if firstAncestor != nil {
					firstAncestor.Children = append(firstAncestor.Children, newItem)
					newItem.Parent = firstAncestor
				} else {
					stack = append(stack, newItem)
				}
			}
		} else {
			stack = append(stack, newItem)
		}

		lastInserted = newItem
	}

	err := ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch n.Kind() {
		case ast.KindHeading:
			heading, ok := n.(*ast.Heading)
			if !ok {
				return ast.WalkStop, errors.Errorf("unexpected node type '%T'", n)
			}

			appendHeading(heading)

		case directive.KindDirective:
			includedNode, exists := include.IncludedNode(n)
			if !exists {
				return ast.WalkContinue, nil
			}

			includedSource, exists := include.IncludedSource(n)
			if !exists {
				return ast.WalkContinue, nil
			}

			includedReader := text.NewReader(includedSource)

			includedToc, err := buildTree(includedNode, includedReader, minLevel, maxLevel)
			if err != nil {
				return ast.WalkStop, errors.WithStack(err)
			}

			var lastItem *tocItem
			if len(stack) > 0 {
				lastItem = stack[len(stack)-1]
			}

			if lastItem != nil {
				for _, item := range includedToc {
					if item == nil {
						continue
					}

					if lastItem.Level < item.Level {
						lastItem.Children = append(lastItem.Children, item)
					} else {
						stack = append(stack, item)
					}
				}
			} else {
				stack = append(stack, includedToc...)
			}

		default:
			return ast.WalkContinue, nil
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return stack, nil
}

func getNodeMinLevelAttribute(node ast.Node) (int, error) {
	minLevelAttrValue, exists := node.AttributeString(attrNameMinLevel)
	if !exists {
		return 1, nil
	}

	rawMinLevel, ok := minLevelAttrValue.(string)
	if !ok {
		return 0, errors.Errorf("unexpected value type '%T' for '%s' attribute", minLevelAttrValue, attrNameMinLevel)
	}

	minLevel, err := strconv.ParseInt(rawMinLevel, 10, 32)
	if err != nil {
		return 0, errors.Wrapf(err, "could not parse value '%v' for '%s' attribute", rawMinLevel, attrNameMinLevel)
	}

	return int(minLevel), nil
}

const attrNameMaxLevel = "maxLevel"

func getNodeMaxLevelAttribute(node ast.Node) (int, error) {
	maxLevelAttrValue, exists := node.AttributeString(attrNameMaxLevel)
	if !exists {
		return math.MaxInt, nil
	}

	rawMaxLevel, ok := maxLevelAttrValue.(string)
	if !ok {
		return 0, errors.Errorf("unexpected value type '%T' for '%s' attribute", maxLevelAttrValue, attrNameMaxLevel)
	}

	maxLevel, err := strconv.ParseInt(rawMaxLevel, 10, 32)
	if err != nil {
		return 0, errors.Wrapf(err, "could not parse value '%v' for '%s' attribute", rawMaxLevel, attrNameMaxLevel)
	}

	return int(maxLevel), nil
}
