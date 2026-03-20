package selector

import (
	"strings"

	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
)

func nodeMatches(node ast.Node, sel simpleSelector, source []byte) bool {
	if sel.typeStr != "" && !matchesType(node, sel.typeStr) {
		return false
	}
	if sel.id != "" {
		idAttr, exists := node.AttributeString("id")
		if !exists {
			return false
		}
		var idStr string
		switch v := idAttr.(type) {
		case []byte:
			idStr = string(v)
		case string:
			idStr = v
		default:
			return false
		}
		if idStr != sel.id {
			return false
		}
	}
	if sel.attrName != "" && !matchesAttr(node, sel.attrName, sel.attrVal, source) {
		return false
	}
	if sel.contains != "" {
		text := strings.TrimSpace(string(node.Text(source)))
		if !globMatch(sel.contains, text) {
			return false
		}
	}
	return true
}

func matchesType(node ast.Node, typeStr string) bool {
	switch typeStr {
	case "h1", "h2", "h3", "h4", "h5", "h6":
		heading, ok := node.(*ast.Heading)
		if !ok {
			return false
		}
		return heading.Level == int(typeStr[1]-'0')
	case "p":
		return node.Kind() == ast.KindParagraph
	case "blockquote":
		return node.Kind() == ast.KindBlockquote
	case "code":
		return node.Kind() == ast.KindFencedCodeBlock
	case "ul":
		list, ok := node.(*ast.List)
		return ok && !list.IsOrdered()
	case "ol":
		list, ok := node.(*ast.List)
		return ok && list.IsOrdered()
	case "table":
		return node.Kind() == extast.KindTable
	}
	return false
}

func matchesAttr(node ast.Node, name, value string, source []byte) bool {
	switch name {
	case "lang":
		fcb, ok := node.(*ast.FencedCodeBlock)
		if !ok {
			return false
		}
		return string(fcb.Language(source)) == value
	default:
		attr, exists := node.AttributeString(name)
		if !exists {
			return false
		}
		switch v := attr.(type) {
		case []byte:
			return string(v) == value
		case string:
			return v == value
		}
	}
	return false
}

// globMatch matches text against a pattern where '*' matches any sequence of characters.
func globMatch(pattern, text string) bool {
	return matchGlob(pattern, text)
}

func matchGlob(pattern, text string) bool {
	for len(pattern) > 0 {
		if pattern[0] == '*' {
			pattern = pattern[1:]
			if len(pattern) == 0 {
				return true
			}
			for i := 0; i <= len(text); i++ {
				if matchGlob(pattern, text[i:]) {
					return true
				}
			}
			return false
		}
		if len(text) == 0 || pattern[0] != text[0] {
			return false
		}
		pattern = pattern[1:]
		text = text[1:]
	}
	return len(text) == 0
}
