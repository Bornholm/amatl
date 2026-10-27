package mcp

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const containerDocument = `# Title

## Offer

:::columns{count=2}
:::card
### Libre
Open formats.
:::
:::card
### Souverain
- Hosted in France
- Under your control
:::
:::

## Contact

Write to us.
`

func newTestServer(t *testing.T, name string, content string) *mcpServer {
	t.Helper()

	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatalf("%+v", err)
	}

	root, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	t.Cleanup(func() { root.Close() })

	return &mcpServer{root: root}
}

func TestReadSectionWithContainer(t *testing.T) {
	srv := newTestServer(t, "doc.md", containerDocument)

	_, out, err := srv.handleReadSection(context.Background(), nil, selectorInput{
		File:     "doc.md",
		Selector: "h2#offer",
	})
	if err != nil {
		t.Fatalf("%+v", err)
	}

	for _, expected := range []string{":::columns{count=\"2\"}", ":::card", "### Libre", "- Under your control"} {
		if !strings.Contains(out.Markdown, expected) {
			t.Errorf("read_section: expected %q in:\n%s", expected, out.Markdown)
		}
	}

	if strings.Contains(out.Markdown, "Contact") {
		t.Errorf("read_section: unexpected next section in:\n%s", out.Markdown)
	}
}

func TestListSectionsWithContainer(t *testing.T) {
	srv := newTestServer(t, "doc.md", containerDocument)

	_, out, err := srv.handleListSections(context.Background(), nil, fileInput{File: "doc.md"})
	if err != nil {
		t.Fatalf("%+v", err)
	}

	var titles []string
	for _, s := range out.Sections {
		titles = append(titles, s.Title)
	}

	if e, g := "Title,Offer,Libre,Souverain,Contact", strings.Join(titles, ","); e != g {
		t.Errorf("list_sections: expected '%s', got '%s'", e, g)
	}
}

func TestUpdateSectionWithContainer(t *testing.T) {
	srv := newTestServer(t, "doc.md", containerDocument)

	_, _, err := srv.handleUpdateSection(context.Background(), nil, updateInput{
		File:     "doc.md",
		Selector: "h2#offer",
		Content:  "## Offer\n\nReplaced.\n",
	})
	if err != nil {
		t.Fatalf("%+v", err)
	}

	data, err := srv.root.ReadFile("doc.md")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	expected := "# Title\n\n## Offer\n\nReplaced.\n\n## Contact\n\nWrite to us.\n"
	if g := string(data); g != expected {
		t.Errorf("update_section: expected:\n%s\ngot:\n%s", expected, g)
	}
}

func TestUpdateSectionEnding(t *testing.T) {
	type testCase struct {
		Name     string
		Source   string
		Selector string
		Content  string
		Expected string
	}

	testCases := []testCase{
		{
			Name:     "list",
			Source:   "## First\n\n- one\n- two\n\n## Second\n\nText.\n",
			Expected: "## First\n\nReplaced.\n\n## Second\n\nText.\n",
		},
		{
			Name:     "fenced code block",
			Source:   "## First\n\n```go\nx := 1\n```\n\n## Second\n\nText.\n",
			Expected: "## First\n\nReplaced.\n\n## Second\n\nText.\n",
		},
		{
			Name:     "thematic break",
			Source:   "## First\n\nText\n\n---\n\n## Second\n\nText.\n",
			Expected: "## First\n\nReplaced.\n\n## Second\n\nText.\n",
		},
		{
			Name:     "table",
			Source:   "## First\n\n| a | b |\n|---|---|\n| 1 | 2 |\n\n## Second\n\nText.\n",
			Expected: "## First\n\nReplaced.\n\n## Second\n\nText.\n",
		},
		{
			Name:     "setext heading",
			Source:   "First\n-----\n\n## Second\n\nText.\n",
			Expected: "## First\n\nReplaced.\n\n## Second\n\nText.\n",
		},
		{
			Name:     "unclosed empty container",
			Source:   "## Second\n\nText.\n\n## First\n\n:::card\n",
			Expected: "## Second\n\nText.\n\n## First\n\nReplaced.\n",
		},
		{
			Name:     "front matter",
			Source:   "---\ntitle: Doc\n---\n\n## First\n\nText\n",
			Expected: "---\ntitle: Doc\n---\n\n## First\n\nReplaced.\n",
		},
		{
			Name:     "code block selector",
			Source:   "## First\n\n```go\nx := 1\n```\n\nAfter.\n",
			Selector: `code[lang="go"]`,
			Content:  "```go\ny := 2\n```\n",
			Expected: "## First\n\n```go\ny := 2\n```\n\nAfter.\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			srv := newTestServer(t, "doc.md", tc.Source)

			selector, content := tc.Selector, tc.Content
			if selector == "" {
				selector, content = "h2#first", "## First\n\nReplaced.\n"
			}

			_, _, err := srv.handleUpdateSection(context.Background(), nil, updateInput{
				File:     "doc.md",
				Selector: selector,
				Content:  content,
			})
			if err != nil {
				t.Fatalf("%+v", err)
			}

			data, err := srv.root.ReadFile("doc.md")
			if err != nil {
				t.Fatalf("%+v", err)
			}

			if g := string(data); g != tc.Expected {
				t.Errorf("update_section: expected:\n%q\ngot:\n%q", tc.Expected, g)
			}
		})
	}
}

func TestListSectionsSkipsFrontMatter(t *testing.T) {
	srv := newTestServer(t, "doc.md", "---\ntitle: Doc\n---\n\n## First\n\nText\n")

	_, out, err := srv.handleListSections(context.Background(), nil, fileInput{File: "doc.md"})
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if len(out.Sections) != 1 || out.Sections[0].Title != "First" {
		t.Errorf("list_sections: expected only 'First', got %+v", out.Sections)
	}
}

func TestReadSectionInsideContainer(t *testing.T) {
	srv := newTestServer(t, "doc.md", containerDocument)

	// Only top-level headings start a section
	_, _, err := srv.handleReadSection(context.Background(), nil, selectorInput{
		File:     "doc.md",
		Selector: "h3#libre",
	})
	if err == nil {
		t.Fatal("read_section: expected an error for a heading inside a container")
	}

	_, out, err := srv.handleFindSections(context.Background(), nil, selectorInput{
		File:     "doc.md",
		Selector: "h3#libre",
	})
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if !strings.Contains(out.Markdown, "Open formats.") {
		t.Errorf("find_sections: expected the card text in:\n%s", out.Markdown)
	}
}
