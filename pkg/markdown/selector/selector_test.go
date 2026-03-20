package selector_test

import (
	"testing"

	"github.com/Bornholm/amatl/pkg/markdown/selector"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func parseMarkdown(src string) (ast.Node, []byte) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)
	p := md.Parser()
	p.AddOptions(parser.WithAutoHeadingID())
	source := []byte(src)
	reader := text.NewReader(source)
	doc := p.Parse(reader)
	return doc, source
}

func TestParse(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
		count   int
	}{
		{"h2", false, 1},
		{"#introduction", false, 1},
		{"h2#api", false, 1},
		{`code[lang="go"]`, false, 1},
		{`h2:contains("API")`, false, 1},
		{"h2#intro, h2#conclusion", false, 2},
		{"h2 > h3", false, 1},
		{"h2 ~ table", false, 1},
		{"blockquote p", false, 1},
		{"", true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			sl, err := selector.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && len(sl) != tt.count {
				t.Errorf("Parse(%q) got %d selectors, want %d", tt.input, len(sl), tt.count)
			}
		})
	}
}

func TestSelectorString(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"h2", "h2"},
		{"#intro", "#intro"},
		{"h2#intro", "h2#intro"},
		{`code[lang="go"]`, `code[lang="go"]`},
		{`h2:contains("My title")`, `h2:contains("My title")`},
		{"h2 > h3", "h2 > h3"},
		{"h2 ~ table", "h2 ~ table"},
		{"blockquote p", "blockquote p"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			sl, err := selector.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", tt.input, err)
			}
			if got := sl[0].String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMatchTopLevel_Headings(t *testing.T) {
	src := "# Title\n\n## Introduction\n\nSome paragraph.\n\n## API Reference\n\n```go\nfunc main() {}\n```\n\n## Conclusion\n\nFinal thoughts.\n"
	doc, source := parseMarkdown(src)

	t.Run("match all h2", func(t *testing.T) {
		sl, _ := selector.Parse("h2")
		nodes := sl.MatchTopLevel(doc, source)
		// 3 h2 headings × their section content
		headings := 0
		for _, n := range nodes {
			if n.Kind() == ast.KindHeading {
				headings++
			}
		}
		if headings != 3 {
			t.Errorf("expected 3 h2 headings, got %d", headings)
		}
	})

	t.Run("match h1", func(t *testing.T) {
		sl, _ := selector.Parse("h1")
		nodes := sl.MatchTopLevel(doc, source)
		// h1 section includes heading + all following siblings (no other h1)
		if len(nodes) == 0 {
			t.Error("expected at least 1 node for h1, got 0")
		}
		if nodes[0].Kind() != ast.KindHeading {
			t.Error("first node should be the h1 heading")
		}
		h := nodes[0].(*ast.Heading)
		if h.Level != 1 {
			t.Errorf("expected level 1 heading, got level %d", h.Level)
		}
	})

	t.Run("match by id", func(t *testing.T) {
		sl, _ := selector.Parse("#introduction")
		nodes := sl.MatchTopLevel(doc, source)
		if len(nodes) == 0 {
			t.Error("expected at least 1 node for #introduction, got 0")
		}
	})
}

func TestMatchTopLevel_Code(t *testing.T) {
	src := "# Title\n\n```go\nfunc main() {}\n```\n\n```python\nprint('hello')\n```\n"
	doc, source := parseMarkdown(src)

	sl, _ := selector.Parse(`code[lang="go"]`)
	nodes := sl.MatchTopLevel(doc, source)
	if len(nodes) != 1 {
		t.Errorf(`expected 1 code[lang="go"] node, got %d`, len(nodes))
	}
}

func TestMatchTopLevel_Contains(t *testing.T) {
	src := "# My Title\n\n## API Reference\n\nSome content.\n\n## Other Section\n\nOther content.\n"
	doc, source := parseMarkdown(src)

	sl, _ := selector.Parse(`h2:contains("API *")`)
	nodes := sl.MatchTopLevel(doc, source)
	headings := 0
	for _, n := range nodes {
		if n.Kind() == ast.KindHeading {
			headings++
		}
	}
	if headings != 1 {
		t.Errorf("expected 1 heading for :contains(\"API *\"), got %d", headings)
	}
}

func TestMatchTopLevel_MultiSelector(t *testing.T) {
	src := "# Title\n\n## Introduction\n\nParagraph.\n\n## Conclusion\n\nEnd.\n"
	doc, source := parseMarkdown(src)

	sl, _ := selector.Parse("#introduction, #conclusion")
	nodes := sl.MatchTopLevel(doc, source)
	headings := 0
	for _, n := range nodes {
		if n.Kind() == ast.KindHeading {
			headings++
		}
	}
	if headings != 2 {
		t.Errorf("expected 2 headings, got %d", headings)
	}
}

func TestSectionNodes(t *testing.T) {
	src := "## Section A\n\nParagraph A.\n\n## Section B\n\nParagraph B.\n"
	doc, _ := parseMarkdown(src)

	// Find first h2
	var h2 *ast.Heading
	for n := doc.FirstChild(); n != nil; n = n.NextSibling() {
		if heading, ok := n.(*ast.Heading); ok && heading.Level == 2 {
			h2 = heading
			break
		}
	}
	if h2 == nil {
		t.Fatal("no h2 found")
	}

	nodes := selector.SectionNodes(h2)
	// Should be: [h2 "Section A", Paragraph "Paragraph A."]
	if len(nodes) != 2 {
		t.Errorf("expected 2 section nodes, got %d", len(nodes))
	}
	if nodes[0].Kind() != ast.KindHeading {
		t.Error("first node should be a heading")
	}
}

func TestGlobMatch(t *testing.T) {
	src := "## API Reference\n\nContent.\n"
	doc, source := parseMarkdown(src)

	tests := []struct {
		pattern string
		wantMatch bool
	}{
		{"API Reference", true},
		{"API *", true},
		{"*Reference", true},
		{"*API*", true},
		{"Other", false},
	}

	var h2 *ast.Heading
	for n := doc.FirstChild(); n != nil; n = n.NextSibling() {
		if heading, ok := n.(*ast.Heading); ok {
			h2 = heading
			break
		}
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			sl, err := selector.Parse(`h2:contains("` + tt.pattern + `")`)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			got := sl[0].Matches(h2, source)
			if got != tt.wantMatch {
				t.Errorf("Matches(%q) = %v, want %v", tt.pattern, got, tt.wantMatch)
			}
		})
	}
}
