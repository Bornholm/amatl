package render

import (
	"github.com/Bornholm/amatl/pkg/markdown/directive"
	"github.com/Bornholm/amatl/pkg/markdown/directive/include"
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"

	highlighting "github.com/yuin/goldmark-highlighting/v2"
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

	return render
}

func newHTMLRenderer() renderer.Renderer {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.Highlighting,
		),
	)

	return markdown.Renderer()
}
