package selector

import (
	"strings"

	"github.com/yuin/goldmark/ast"
)

type combinatorKind int

const (
	combinatorDescendant combinatorKind = iota // A B (space)
	combinatorChild                            // A > B
	combinatorSibling                          // A ~ B
)

type simpleSelector struct {
	typeStr  string // "h1", "p", etc.
	id       string // "" if not specified
	attrName string // "" if not specified
	attrVal  string
	contains string // "" if not specified
}

func (ss simpleSelector) String() string {
	var sb strings.Builder
	sb.WriteString(ss.typeStr)
	if ss.id != "" {
		sb.WriteByte('#')
		sb.WriteString(ss.id)
	}
	if ss.attrName != "" {
		sb.WriteString(`[`)
		sb.WriteString(ss.attrName)
		sb.WriteString(`="`)
		sb.WriteString(ss.attrVal)
		sb.WriteString(`"]`)
	}
	if ss.contains != "" {
		sb.WriteString(`:contains("`)
		sb.WriteString(ss.contains)
		sb.WriteString(`")`)
	}
	return sb.String()
}

// compoundSelector is a selector with optional left-hand combinator chain.
type compoundSelector struct {
	left       *compoundSelector
	combinator combinatorKind
	right      simpleSelector
}

// Matches checks if the given node matches this compound selector.
// Evaluation is right-to-left: check right matches node, then check left with combinator.
func (cs *compoundSelector) Matches(node ast.Node, source []byte) bool {
	if !nodeMatches(node, cs.right, source) {
		return false
	}
	if cs.left == nil {
		return true
	}
	switch cs.combinator {
	case combinatorDescendant:
		for p := node.Parent(); p != nil; p = p.Parent() {
			if cs.left.Matches(p, source) {
				return true
			}
		}
		return false
	case combinatorChild:
		p := node.Parent()
		if p == nil {
			return false
		}
		return cs.left.Matches(p, source)
	case combinatorSibling:
		for sibling := node.PreviousSibling(); sibling != nil; sibling = sibling.PreviousSibling() {
			if cs.left.Matches(sibling, source) {
				return true
			}
		}
		return false
	}
	return false
}

// String returns the string representation of this compound selector.
func (cs *compoundSelector) String() string {
	if cs.left == nil {
		return cs.right.String()
	}
	switch cs.combinator {
	case combinatorChild:
		return cs.left.String() + " > " + cs.right.String()
	case combinatorSibling:
		return cs.left.String() + " ~ " + cs.right.String()
	default: // combinatorDescendant
		return cs.left.String() + " " + cs.right.String()
	}
}

var _ Selector = &compoundSelector{}
