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
	"github.com/pkg/errors"
)

// TestContainerDirectivePipeline checks that container directives survive
// the intermediate Markdown stage and integrate with other directives,
// as the html and pdf commands chain them.
func TestContainerDirectivePipeline(t *testing.T) {
	html := renderContainerTestData(t, "testdata/container")

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

	if e, g := expected, html; e != g {
		t.Errorf("unexpected html:\n%s", diff.LineDiff(e, g))
	}
}

// TestContainerBackgroundPipeline checks that the bg attribute is resolved
// against the file declaring it, included files too, and embedded.
func TestContainerBackgroundPipeline(t *testing.T) {
	html := renderContainerTestData(t, "testdata/container-bg")

	dataURL := `url(&quot;data:image/svg+xml;base64,`

	expected := []string{
		`<section class="amatl-section hero" style="--amatl-bg: ` + dataURL,
		`<div class="amatl-card" style="color: red; --amatl-bg: ` + dataURL,
	}

	for _, e := range expected {
		if !strings.Contains(html, e) {
			t.Errorf("expected %q in html:\n%s", e, html)
		}
	}

	if strings.Contains(html, "data-bg") {
		t.Errorf("unexpected data-bg attribute in html:\n%s", html)
	}
}

// TestContainerBackgroundMissing checks that a bg pointing to a missing
// file fails the render with an error instead of a panic
func TestContainerBackgroundMissing(t *testing.T) {
	_, err := transformContainerTestData(t, "testdata/container-bg-missing")
	if err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(err.Error(), "could not resolve resource './missing.png'") {
		t.Errorf("unexpected error: %+v", err)
	}
}

// renderContainerTestData renders the main.md file of the given directory
// with its layout.html, as the html command does
func renderContainerTestData(t *testing.T, dir string) string {
	t.Helper()

	html, err := transformContainerTestData(t, dir)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	return html
}

func transformContainerTestData(t *testing.T, dir string) (string, error) {
	t.Helper()

	sourcePath, err := resolver.Path(dir + "/main.md").Abs()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	layoutPath, err := resolver.Path(dir + "/layout.html").Abs()
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
		return "", errors.WithStack(err)
	}

	return strings.TrimSpace(string(payload.GetData())), nil
}
