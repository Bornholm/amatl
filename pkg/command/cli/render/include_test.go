package render

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Bornholm/amatl/pkg/pipeline"
	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/andreyvit/diff"
)

// TestIncludeSameRelativeURL checks that two files of different directories
// including the same relative URL each get their own file, with links
// rewritten against its directory
func TestIncludeSameRelativeURL(t *testing.T) {
	sourcePath, err := resolver.Path("testdata/include-same-name/main.md").Abs()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	baseDir, err := sourcePath.Dir().Abs()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	source, err := os.ReadFile(sourcePath.String())
	if err != nil {
		t.Fatalf("%+v", err)
	}

	ctx := resolver.WithWorkDir(context.Background(), baseDir)

	transformer := pipeline.Pipeline(
		MarkdownMiddleware(WithSourcePath(sourcePath)),
	)

	payload := pipeline.NewPayload(source)

	if err := transformer.Transform(ctx, payload); err != nil {
		t.Fatalf("%+v", err)
	}

	// Included paths use the separator of the platform
	markdown := strings.ReplaceAll(
		filepath.ToSlash(string(payload.GetData())),
		filepath.ToSlash(baseDir.String()),
		"<dir>",
	)

	for _, expected := range []string{
		"## Doc a\n\nPart of a\n\n![a](<dir>/a/img.svg)",
		"## Doc c\n\nPart of c\n\n![c](<dir>/c/img.svg)",
	} {
		if !strings.Contains(markdown, expected) {
			t.Errorf("expected %q in markdown:\n%s", expected, diff.LineDiff(expected, markdown))
		}
	}
}

// TestIncludeSameRelativeURLHTML checks the same case through the HTML stage,
// where each image must be embedded from its own directory
func TestIncludeSameRelativeURLHTML(t *testing.T) {
	html := renderContainerTestData(t, "testdata/include-same-name")

	for _, dir := range []string{"a", "c"} {
		data, err := os.ReadFile(filepath.Join("testdata/include-same-name", dir, "img.svg"))
		if err != nil {
			t.Fatalf("%+v", err)
		}

		expected := `src="data:image/svg+xml;base64,` + base64.StdEncoding.EncodeToString(data) + `" alt="` + dir + `"`
		if !strings.Contains(html, expected) {
			t.Errorf("expected the image of '%s' in html:\n%s", dir, html)
		}
	}
}

// TestIncludeCycle checks that a document including itself through another
// one fails with an error instead of recursing until the stack overflows
func TestIncludeCycle(t *testing.T) {
	_, err := transformContainerTestData(t, "testdata/include-cycle")
	if err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(err.Error(), "include cycle:") || !strings.Contains(err.Error(), "b.md -> ") {
		t.Errorf("unexpected error: %+v", err)
	}
}
