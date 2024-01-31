package render

import (
	"bytes"
	"context"

	"github.com/Bornholm/amatl/pkg/html"
	"github.com/Bornholm/amatl/pkg/pipeline"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/text"
)

type MarkdownTransformerOptions struct {
	BaseDir string
	WithToc bool
}

type MarkdownTransformerOptionFunc func(opts *MarkdownTransformerOptions)

func NewMarkdownTransformerOptions(funcs ...MarkdownTransformerOptionFunc) *MarkdownTransformerOptions {
	opts := &MarkdownTransformerOptions{}
	for _, fn := range funcs {
		fn(opts)
	}
	return opts
}

func WithBaseDir(baseDir string) MarkdownTransformerOptionFunc {
	return func(opts *MarkdownTransformerOptions) {
		opts.BaseDir = baseDir
	}
}

func WithToc(enabled bool) MarkdownTransformerOptionFunc {
	return func(opts *MarkdownTransformerOptions) {
		opts.WithToc = enabled
	}
}

func MarkdownTransformer(funcs ...MarkdownTransformerOptionFunc) pipeline.Transformer {
	opts := NewMarkdownTransformerOptions(funcs...)
	return pipeline.NewTransformer(func(ctx context.Context, input []byte) ([]byte, error) {
		reader := text.NewReader(input)

		parse := newParser(opts.BaseDir, opts.WithToc)
		render := newMarkdownRenderer()

		document := parse.Parse(reader)

		var buf bytes.Buffer
		if err := render.Render(&buf, input, document); err != nil {
			return nil, errors.WithStack(err)
		}

		return buf.Bytes(), nil
	})
}

type HTMLTransformerOptions struct {
	*MarkdownTransformerOptions
	LayoutURL  string
	LayoutVars map[string]any
}

type HTMLTransformerOptionFunc func(opts *HTMLTransformerOptions)

func WithMarkdownTransformerOptions(funcs ...MarkdownTransformerOptionFunc) HTMLTransformerOptionFunc {
	return func(opts *HTMLTransformerOptions) {
		opts.MarkdownTransformerOptions = NewMarkdownTransformerOptions(funcs...)
	}
}

func NewHTMLTransformerOptions(funcs ...HTMLTransformerOptionFunc) *HTMLTransformerOptions {
	opts := &HTMLTransformerOptions{
		MarkdownTransformerOptions: NewMarkdownTransformerOptions(),
		LayoutURL:                  html.DefaultLayoutURL,
		LayoutVars:                 make(map[string]any),
	}
	for _, fn := range funcs {
		fn(opts)
	}
	return opts
}

func WithLayoutURL(layoutURL string) HTMLTransformerOptionFunc {
	return func(opts *HTMLTransformerOptions) {
		opts.LayoutURL = layoutURL
	}
}

func WithLayoutVars(vars map[string]any) HTMLTransformerOptionFunc {
	return func(opts *HTMLTransformerOptions) {
		opts.LayoutVars = vars
	}
}

func HTMLTransformer(funcs ...HTMLTransformerOptionFunc) pipeline.Transformer {
	opts := NewHTMLTransformerOptions(funcs...)
	return pipeline.NewTransformer(func(ctx context.Context, input []byte) ([]byte, error) {
		reader := text.NewReader(input)

		parser := newParser(opts.BaseDir, opts.WithToc)
		document := parser.Parse(reader)

		render := newHTMLRenderer()

		var body bytes.Buffer
		if err := render.Render(&body, input, document); err != nil {
			return nil, errors.WithStack(err)
		}

		var layout bytes.Buffer

		err := html.Layout(
			&layout, body.Bytes(),
			html.WithLayout(opts.LayoutURL),
			html.WithVars(opts.LayoutVars),
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return layout.Bytes(), nil
	})
}
