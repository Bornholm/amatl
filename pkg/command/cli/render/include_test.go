package render

import (
	"context"
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
		"## Doc a\n\nPart of a\n\n![a](<dir>/a/img.png)",
		"## Doc c\n\nPart of c\n\n![c](<dir>/c/img.png)",
	} {
		if !strings.Contains(markdown, expected) {
			t.Errorf("expected %q in markdown:\n%s", expected, diff.LineDiff(expected, markdown))
		}
	}
}
