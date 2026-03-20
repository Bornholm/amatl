package selector

import "github.com/yuin/goldmark/ast"

// Selector represents a CSS-like selector for Markdown AST nodes.
type Selector interface {
	Matches(node ast.Node, source []byte) bool
	String() string
}

// SelectorList is a comma-separated list of selectors.
type SelectorList []Selector

// Parse parses a CSS-like selector string into a SelectorList.
func Parse(input string) (SelectorList, error) {
	p := newParser(input)
	return p.parseSelectorList()
}

// MatchTopLevel returns all direct children of root that match the selector list.
// For heading matches, the full section is expanded (heading + subsequent siblings
// until a heading of the same or higher level).
func (sl SelectorList) MatchTopLevel(root ast.Node, source []byte) []ast.Node {
	var result []ast.Node
	seen := make(map[ast.Node]bool)

	for node := root.FirstChild(); node != nil; node = node.NextSibling() {
		if seen[node] {
			continue
		}
		for _, sel := range sl {
			if sel.Matches(node, source) {
				if heading, ok := node.(*ast.Heading); ok {
					for _, n := range SectionNodes(heading) {
						if !seen[n] {
							seen[n] = true
							result = append(result, n)
						}
					}
				} else {
					seen[node] = true
					result = append(result, node)
				}
				break
			}
		}
	}

	return result
}

// FindAll returns all nodes in root (walking the full tree) that match the selector list.
// For heading matches, the full section is expanded.
func (sl SelectorList) FindAll(root ast.Node, source []byte) []ast.Node {
	var result []ast.Node
	seen := make(map[ast.Node]bool)

	_ = ast.Walk(root, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || node == root {
			return ast.WalkContinue, nil
		}
		if seen[node] {
			return ast.WalkSkipChildren, nil
		}
		for _, sel := range sl {
			if sel.Matches(node, source) {
				if heading, ok := node.(*ast.Heading); ok {
					for _, n := range SectionNodes(heading) {
						if !seen[n] {
							seen[n] = true
							result = append(result, n)
						}
					}
					return ast.WalkSkipChildren, nil
				}
				seen[node] = true
				result = append(result, node)
				return ast.WalkSkipChildren, nil
			}
		}
		return ast.WalkContinue, nil
	})

	return result
}
