package amatl

import (
	"regexp"
	"strings"
	"testing"
)

var (
	// cssRule matches a CSS rule, capturing its selector list and its body
	cssRule = regexp.MustCompile(`([^{}]+)\{([^}]*)\}`)
	// cssComment matches a CSS comment
	cssComment = regexp.MustCompile(`(?s)/\*.*?\*/`)
)

// cssRules returns the body of each rule of the given CSS whose selector
// list holds the given selector, ignoring whitespace
func cssRules(css string, selector string) []string {
	normalize := func(s string) string {
		return strings.Join(strings.Fields(s), " ")
	}

	var bodies []string

	css = cssComment.ReplaceAllString(css, "")

	for _, match := range cssRule.FindAllStringSubmatch(css, -1) {
		for _, s := range strings.Split(match[1], ",") {
			if normalize(s) == selector {
				bodies = append(bodies, match[2])
				break
			}
		}
	}

	return bodies
}

// TestPresentationNestedRules checks that the presentation layout does not
// hide or break the page on a bare hr: htmlSplit drops the top-level ones,
// so such a rule would only hit a hr nested in a container or a blockquote.
func TestPresentationNestedRules(t *testing.T) {
	data, err := templateFs.ReadFile("templates/presentation.html")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	for _, body := range cssRules(string(data), "hr") {
		for _, property := range []string{"page-break-after", "page-break-before", "break-after", "break-before", "opacity", "display", "visibility"} {
			if strings.Contains(body, property) {
				t.Errorf("presentation.html: unexpected '%s' in a hr rule: %q", property, strings.TrimSpace(body))
			}
		}
	}
}

// TestComponentsGridRules checks that a hr spans the whole row of columns
// and grid containers instead of taking a cell
func TestComponentsGridRules(t *testing.T) {
	data, err := templateFs.ReadFile("assets/components.css")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	gridColumn := regexp.MustCompile(`grid-column\s*:\s*1\s*/\s*-1`)

	for _, selector := range []string{".amatl-columns > hr", ".amatl-grid > hr"} {
		found := false
		for _, body := range cssRules(string(data), selector) {
			if gridColumn.MatchString(body) {
				found = true
			}
		}

		if !found {
			t.Errorf("components.css: expected '%s' to span the whole row", selector)
		}
	}
}

func TestCSSRules(t *testing.T) {
	css := `section.slide > hr { opacity: 0 } hr, .slide { display: none } /* c */ .a  >  hr { x: 1 }`

	if g := cssRules(css, "hr"); len(g) != 1 || !strings.Contains(g[0], "display") {
		t.Errorf("cssRules: expected the grouped bare hr rule only, got %q", g)
	}

	if g := cssRules(css, ".a > hr"); len(g) != 1 {
		t.Errorf("cssRules: expected one rule, got %q", g)
	}
}
