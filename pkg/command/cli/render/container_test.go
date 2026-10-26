package render

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/Bornholm/amatl/pkg/markdown/directive/attrs"
	"github.com/Bornholm/amatl/pkg/markdown/directive/toc"
	"github.com/Bornholm/amatl/pkg/pipeline"
	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/andreyvit/diff"
)

// TestContainerDirectivePipeline checks that container directives survive
// the intermediate Markdown stage and integrate with other directives,
// as the html and pdf commands chain them.
func TestContainerDirectivePipeline(t *testing.T) {
	sourcePath, err := resolver.Path("testdata/container/main.md").Abs()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	layoutPath, err := resolver.Path("testdata/container/layout.html").Abs()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	source, err := os.ReadFile(sourcePath.String())
	if err != nil {
		t.Fatalf("%+v", err)
	}

	baseDir, err := sourcePath.Dir().Abs()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	ctx := resolver.WithWorkDir(context.Background(), baseDir)

	transformer := pipeline.Pipeline(
		MarkdownMiddleware(
			WithSourcePath(sourcePath),
			WithIgnoredDirectives(toc.Type, attrs.Type),
		),
		HTMLMiddleware(
			WithMarkdownTransformerOptions(
				WithSourcePath(sourcePath),
			),
			WithLayoutURL(layoutPath.String()),
		),
	)

	payload := pipeline.NewPayload(source)

	if err := transformer.Transform(ctx, payload); err != nil {
		t.Fatalf("%+v", err)
	}

	// The empty paragraph is left by the :attrs{} directive line,
	// independently of containers
	expected := `<h1 id="brochure">Brochure</h1>
<p></p>
<div class="amatl-columns highlight" data-count="2">
<div class="amatl-card">
<p>Included <strong>content</strong></p>
</div>
<div class="amatl-card">
<p>Second</p>
</div>
</div>`

	if e, g := expected, strings.TrimSpace(string(payload.GetData())); e != g {
		t.Errorf("unexpected html:\n%s", diff.LineDiff(e, g))
	}
}
