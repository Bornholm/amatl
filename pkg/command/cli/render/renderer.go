package render

import (
	"github.com/Bornholm/amatl/pkg/markdown/directive"
	"github.com/Bornholm/amatl/pkg/markdown/directive/include"
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"

	highlighting "github.com/yuin/goldmark-highlighting/v2"

	mermaidRenderer "github.com/Bornholm/amatl/pkg/markdown/renderer/markdown/mermaid"
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown/node"
	"go.abhg.dev/goldmark/mermaid"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
)

func newMarkdownRenderer() renderer.Renderer {
	render := markdown.NewRenderer()

	render.AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(
				directive.NewRenderer(
					directive.WithRenderer(
						include.Type,
						&include.NodeRenderer{
							Cache:    cache,
							Renderer: render,
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
		markdown.WithNodeRenderer(
			mermaid.Kind,
			&mermaidRenderer.BlockNodeRenderer{},
		),
		markdown.WithNodeRenderer(
			mermaid.ScriptKind,
			&mermaidRenderer.ScriptBlockNodeRenderer{},
		),
	)

	return render
}

func newHTMLRenderer() renderer.Renderer {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("nord"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(false),
				),
			),
			&mermaid.Extender{
				RenderMode: mermaid.RenderModeClient,
			},
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	return markdown.Renderer()
}
