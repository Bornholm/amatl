package amatl

import (
	"regexp"
	"strings"
	"testing"
)

// unscopedHRRule matches a CSS rule whose selector is a bare hr
var unscopedHRRule = regexp.MustCompile(`(?:^|[\s,}])hr\s*\{([^}]*)\}`)

// TestPresentationNestedRules checks that the presentation layout does not
// hide or break the page on a hr: htmlSplit drops the top-level ones, so
// such a rule would only hit a hr nested in a container or a blockquote.
func TestPresentationNestedRules(t *testing.T) {
	data, err := templateFs.ReadFile("templates/presentation.html")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	for _, match := range unscopedHRRule.FindAllStringSubmatch(string(data), -1) {
		for _, property := range []string{"page-break", "break-after", "break-before", "opacity", "display", "visibility"} {
			if strings.Contains(match[1], property) {
				t.Errorf("presentation.html: unexpected '%s' in rule %q", property, strings.TrimSpace(match[0]))
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

	rule := regexp.MustCompile(`\.amatl-columns > hr,\s*\.amatl-grid > hr\s*\{\s*grid-column: 1 / -1;\s*\}`)
	if !rule.Match(data) {
		t.Error("components.css: expected hr to span the whole row of columns and grid")
	}
}
