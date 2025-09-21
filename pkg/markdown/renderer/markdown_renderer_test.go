package markdown

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Bornholm/amatl/pkg/markdown/directive"
	"github.com/Bornholm/amatl/pkg/markdown/directive/include"
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown/node"
	"github.com/Bornholm/amatl/pkg/transform"
	"github.com/andreyvit/diff"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func TestMarkdonwRenderer(t *testing.T) {
	files, err := filepath.Glob("testdata/renderer/markdown/*.md")
	if err != nil {
		t.Fatalf("%+v", errors.WithStack(err))
	}

	for _, f := range files {
		filename := filepath.Base(f)
		t.Run(filename, func(t *testing.T) {

			file, err := os.Open(f)
			if err != nil {
				t.Fatalf("%+v", errors.WithStack(err))
			}

			transformed := transform.NewNewlineReader(file)

			data, err := io.ReadAll(transformed)
			if err != nil {
				t.Fatalf("%+v", errors.WithStack(err))
			}

			reader := text.NewReader(data)

			gm := goldmark.New(
				goldmark.WithExtensions(
					extension.GFM,
				),
			)

			parser := gm.Parser()

			root := parser.Parse(reader)

			cache := include.NewSourceCache()

			markdownRenderer := markdown.NewRenderer()
			markdownRenderer.AddOptions(
				markdown.WithNodeRenderers(node.Renderers()),
				renderer.WithNodeRenderers(
					util.Prioritized(
						directive.NewRenderer(
							directive.WithRenderer(
								include.Type,
								&include.NodeRenderer{
									Cache:    cache,
									Renderer: markdownRenderer,
								},
							),
						), 0,
					),
				),
				markdown.WithNodeRenderers(node.Renderers()),
				markdown.WithNodeRenderer(
					directive.KindDirective,
					directive.NewMarkdownNodeRenderer(
						directive.WithMarkdownDirectiveRenderer(
							include.Type,
							&include.MarkdownRenderer{
								Cache: cache,
							},
						),
					),
				),
			)

			var buff bytes.Buffer

			if err := markdownRenderer.Render(&buff, data, root); err != nil {
				t.Fatalf("%+v", errors.WithStack(err))
			}

			if e, g := strings.TrimSpace(string(data)), strings.TrimSpace(buff.String()); e != g {
				t.Errorf("Result not as expected:\n%v", diff.LineDiff(e, g))
			}
		})
	}
}
