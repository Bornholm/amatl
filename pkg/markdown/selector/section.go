package selector

import "github.com/yuin/goldmark/ast"

// SectionNodes returns the heading and all subsequent sibling nodes that belong
// to its section (up to but not including the next heading of the same or higher level).
func SectionNodes(heading *ast.Heading) []ast.Node {
	result := []ast.Node{heading}
	for node := heading.NextSibling(); node != nil; node = node.NextSibling() {
		if h, ok := node.(*ast.Heading); ok && h.Level <= heading.Level {
			break
		}
		result = append(result, node)
	}
	return result
}
