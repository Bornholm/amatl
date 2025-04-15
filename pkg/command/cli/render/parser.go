package render

import (
	"net/url"

	"github.com/Bornholm/amatl/pkg/markdown/dataurl"
	"github.com/Bornholm/amatl/pkg/markdown/directive"
	"github.com/Bornholm/amatl/pkg/markdown/directive/attrs"
	"github.com/Bornholm/amatl/pkg/markdown/directive/include"
	"github.com/Bornholm/amatl/pkg/markdown/directive/toc"
	"github.com/Bornholm/amatl/pkg/markdown/linkrewriter"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/frontmatter"
	"go.abhg.dev/goldmark/mermaid"
)

var (
	cache = include.NewSourceCache()
)

type ParserOptions struct {
	EmbedLinkedResources bool
	LinkReplacements     map[string]string
}

func newParser(SourceURL *url.URL, opts ParserOptions) parser.Parser {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			&mermaid.Extender{
				RenderMode: mermaid.RenderModeClient,
			},
			&frontmatter.Extender{
				Mode: frontmatter.SetMetadata,
			},
		),
	)

	parse := markdown.Parser()

	parse.AddOptions(
		parser.WithAutoHeadingID(),
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
					directive.WithTransformer(
						toc.Type,
						&toc.NodeTransformer{},
					),
					directive.WithTransformer(
						attrs.Type,
						&attrs.NodeTransformer{},
					),
				),
				0,
			),
		),
	)

	if opts.EmbedLinkedResources {
		parse.AddOptions(
			parser.WithASTTransformers(
				util.Prioritized(&dataurl.Transformer{}, 990),
			),
		)
	}

	if opts.LinkReplacements != nil {
		parse.AddOptions(
			parser.WithASTTransformers(
				util.Prioritized(linkrewriter.NewTransformer(opts.LinkReplacements), 999),
			),
		)
	}

	return parse
}
