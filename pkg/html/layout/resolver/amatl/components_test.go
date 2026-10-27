package amatl

import (
	"regexp"
	"testing"
)

// TestComponentsBackgroundRule checks that the background rule is scoped to
// containers, which also makes it win over the background shorthand of
// .amatl-card and .amatl-callout
func TestComponentsBackgroundRule(t *testing.T) {
	data, err := templateFs.ReadFile("assets/components.css")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	rule := regexp.MustCompile(`\[class\^="amatl-"\]\[style\*="--amatl-bg:"\],\s*\[class\*=" amatl-"\]\[style\*="--amatl-bg:"\]\s*\{\s*background-image: var\(--amatl-bg\);`)
	if !rule.Match(data) {
		t.Error("components.css: expected a background rule scoped to containers")
	}
}
