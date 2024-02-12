package render

import (
	"net/url"

	"github.com/Bornholm/amatl/pkg/markdown/dataurl"
	"github.com/Bornholm/amatl/pkg/markdown/directive"
	"github.com/Bornholm/amatl/pkg/markdown/directive/include"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/mermaid"
	"go.abhg.dev/goldmark/toc"
)

var (
	cache = include.NewSourceCache()
)

func newParser(SourceURL *url.URL, withToC bool, embedLinkedResources bool) parser.Parser {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			&mermaid.Extender{
				RenderMode: mermaid.RenderModeClient,
			},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	parse := markdown.Parser()

	parse.AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(&directive.InlineParser{}, 0),
		),
		parser.WithASTTransformers(
			util.Prioritized(
				directive.NewTransformer(
					directive.WithTransformer(
						include.Type,
						&include.NodeTransformer{
							SourceURL: SourceURL,
							Cache:     cache,
							Parser:    parse,
						},
					),
				),
				0,
			),
		),
	)

	if withToC {
		parse.AddOptions(
			parser.WithASTTransformers(
				util.Prioritized(&toc.Transformer{
					Title: "Table of content",
				}, 100),
			),
		)
	}

	if embedLinkedResources {
		parse.AddOptions(
			parser.WithASTTransformers(
				util.Prioritized(&dataurl.Transformer{}, 999),
			),
		)
	}

	return parse
}
