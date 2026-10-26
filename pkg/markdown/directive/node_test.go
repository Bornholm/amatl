package directive

import (
	"testing"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func TestParseDirectiveNonStringAttributes(t *testing.T) {
	raw := []byte(`:toc{minLevel=2 enabled=true name="foo"}`)

	node := parseDirective(raw, ast.NewTextSegment(text.NewSegment(0, len(raw))))
	if node == nil {
		t.Fatal("expected directive to be parsed")
	}

	expected := map[string]string{
		"minLevel": "2",
		"enabled":  "true",
		"name":     "foo",
	}

	for name, value := range expected {
		attr, exists := node.AttributeString(name)
		if !exists {
			t.Errorf("attribute '%s' not found", name)
			continue
		}

		if e, g := value, attr.(string); e != g {
			t.Errorf("attribute '%s': expected '%s', got '%s'", name, e, g)
		}
	}
}
