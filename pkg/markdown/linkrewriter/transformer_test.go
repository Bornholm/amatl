package linkrewriter

import (
	"context"
	"testing"

	"github.com/Bornholm/amatl/pkg/pipeline"
	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func TestTransformer(t *testing.T) {
	type testCase struct {
		Name         string
		WorkDir      string
		Replacements map[string]string
		Source       string
		Expected     string
	}

	testCases := []testCase{
		{
			Name:     "relative link",
			WorkDir:  "/home/me/docs",
			Source:   "[a](./guide.md)",
			Expected: "/home/me/docs/guide.md",
		},
		{
			Name:     "encoded space",
			WorkDir:  "/home/me/docs",
			Source:   "[a](./my%20guide.md)",
			Expected: "/home/me/docs/my%20guide.md",
		},
		{
			Name:         "prefix with a space",
			WorkDir:      "/home/me/my project",
			Replacements: map[string]string{"/home/me/my project": "https://example.com/blob/v1"},
			Source:       "[a](<./my guide.md>)",
			Expected:     "https://example.com/blob/v1/my%20guide.md",
		},
		{
			Name:         "file URL workdir",
			WorkDir:      "file:///home/me/my%20project",
			Replacements: map[string]string{"file:///home/me/my project": "https://example.com/blob/v1"},
			Source:       "[a](./guide.md)",
			Expected:     "https://example.com/blob/v1/guide.md",
		},
		{
			Name:         "encoded prefix",
			WorkDir:      "file:///home/me/my%20project",
			Replacements: map[string]string{"file:///home/me/my%20project": "https://example.com/blob/v1"},
			Source:       "[a](./guide.md)",
			Expected:     "https://example.com/blob/v1/guide.md",
		},
		{
			Name:     "absolute URL",
			WorkDir:  "/home/me/docs",
			Source:   "[a](https://example.com/a%20b)",
			Expected: "https://example.com/a%20b",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			parse := goldmark.New().Parser()
			parse.AddOptions(parser.WithASTTransformers(
				util.Prioritized(NewTransformer(tc.Replacements), 0),
			))

			ctx := resolver.WithWorkDir(context.Background(), resolver.Path(tc.WorkDir))
			pc := pipeline.WithContext(ctx, parser.NewContext())

			doc := parse.Parse(text.NewReader([]byte(tc.Source)), parser.WithContext(pc))

			link, ok := doc.FirstChild().FirstChild().(*ast.Link)
			if !ok {
				t.Fatalf("expected *ast.Link, got %T", doc.FirstChild().FirstChild())
			}

			if g := string(link.Destination); g != tc.Expected {
				t.Errorf("expected '%s', got '%s'", tc.Expected, g)
			}
		})
	}
}
