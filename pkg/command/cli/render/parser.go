package render

import (
	"slices"

	"github.com/Bornholm/amatl/pkg/markdown/dataurl"
	"github.com/Bornholm/amatl/pkg/markdown/directive"
	"github.com/Bornholm/amatl/pkg/markdown/directive/attrs"
	"github.com/Bornholm/amatl/pkg/markdown/directive/include"
	"github.com/Bornholm/amatl/pkg/markdown/directive/toc"
	"github.com/Bornholm/amatl/pkg/markdown/linkrewriter"
	"github.com/Bornholm/amatl/pkg/resolver"
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
	IgnoredDirectives    []directive.Type
}

func newParser(sourcePath resolver.Path, opts ParserOptions) parser.Parser {
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

	directiveTransformers := []directive.TransformerOptionFunc{}

	if !isDirectiveIgnored(toc.Type, opts.IgnoredDirectives) {
		directiveTransformers = append(directiveTransformers,
			directive.WithTransformer(
				toc.Type,
				&toc.NodeTransformer{},
			),
		)
	}

	if !isDirectiveIgnored(attrs.Type, opts.IgnoredDirectives) {
		directiveTransformers = append(directiveTransformers,
			directive.WithTransformer(
				attrs.Type,
				&attrs.NodeTransformer{},
			),
		)
	}

	if !isDirectiveIgnored(include.Type, opts.IgnoredDirectives) {
		directiveTransformers = append(directiveTransformers,
			directive.WithTransformer(
				include.Type,
				&include.NodeTransformer{
					SourcePath: sourcePath,
					Cache:      cache,
					Parser:     parse,
				},
			),
		)
	}

	parse.AddOptions(
		parser.WithAutoHeadingID(),
		parser.WithInlineParsers(
			util.Prioritized(&directive.InlineParser{}, 0),
		),
		parser.WithASTTransformers(
			util.Prioritized(
				directive.NewTransformer(directiveTransformers...),
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

func isDirectiveIgnored(dt directive.Type, ignored []directive.Type) bool {
	return slices.ContainsFunc(ignored, func(curr directive.Type) bool {
		return dt == curr
	})
}
