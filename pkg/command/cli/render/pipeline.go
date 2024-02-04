package render

import (
	"bytes"
	"context"
	"sync"

	"github.com/Bornholm/amatl/pkg/html"
	"github.com/Bornholm/amatl/pkg/pipeline"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
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

		parse := newParser(opts.BaseDir, opts.WithToc, false)
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

		parser := newParser(opts.BaseDir, opts.WithToc, true)
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

type PDFTransformerOptions struct {
	MarginTop    float64
	MarginLeft   float64
	MarginRight  float64
	MarginBottom float64
	Scale        float64
}

const (
	DefaultPDFMargin float64 = 1
	DefaultPDFScale  float64 = 1
)

type PDFTransformerOptionFunc func(opts *PDFTransformerOptions)

func NewPDFTransformerOptions(funcs ...PDFTransformerOptionFunc) *PDFTransformerOptions {
	opts := &PDFTransformerOptions{
		MarginTop:    DefaultPDFMargin,
		MarginLeft:   DefaultPDFMargin,
		MarginRight:  DefaultPDFMargin,
		MarginBottom: DefaultPDFMargin,
		Scale:        DefaultPDFScale,
	}
	for _, fn := range funcs {
		fn(opts)
	}
	return opts
}

func WithMarginTop(margin float64) PDFTransformerOptionFunc {
	return func(opts *PDFTransformerOptions) {
		opts.MarginTop = margin
	}
}

func WithMarginLeft(margin float64) PDFTransformerOptionFunc {
	return func(opts *PDFTransformerOptions) {
		opts.MarginLeft = margin
	}
}

func WithMarginBottom(margin float64) PDFTransformerOptionFunc {
	return func(opts *PDFTransformerOptions) {
		opts.MarginBottom = margin
	}
}

func WithMarginRight(margin float64) PDFTransformerOptionFunc {
	return func(opts *PDFTransformerOptions) {
		opts.MarginRight = margin
	}
}

func WithScale(scale float64) PDFTransformerOptionFunc {
	return func(opts *PDFTransformerOptions) {
		opts.Scale = scale
	}
}

func PDFTransformer(funcs ...PDFTransformerOptionFunc) pipeline.Transformer {
	opts := NewPDFTransformerOptions(funcs...)

	return pipeline.NewTransformer(func(ctx context.Context, input []byte) ([]byte, error) {
		var output []byte
		c, _ := chromedp.NewContext(context.Background())
		if err := chromedp.Run(c, printToPDF(input, &output, opts)); err != nil {
			return nil, errors.WithStack(err)
		}

		return output, nil
	})
}

func printToPDF(html []byte, res *[]byte, opts *PDFTransformerOptions) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			lctx, cancel := context.WithCancel(ctx)
			defer cancel()
			var wg sync.WaitGroup
			wg.Add(1)
			chromedp.ListenTarget(lctx, func(ev interface{}) {
				if _, ok := ev.(*page.EventLoadEventFired); ok {
					// It's a good habit to remove the event listener if we don't need it anymore.
					cancel()
					wg.Done()
				}
			})
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			if err := page.SetDocumentContent(frameTree.Frame.ID, string(html)).Do(ctx); err != nil {
				return err
			}
			wg.Wait()
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(false).
				WithDisplayHeaderFooter(false).
				WithPreferCSSPageSize(true).
				WithMarginRight(centimetersToInches(opts.MarginRight)).
				WithMarginTop(centimetersToInches(opts.MarginTop)).
				WithMarginBottom(centimetersToInches(opts.MarginBottom)).
				WithMarginLeft(centimetersToInches(opts.MarginLeft)).
				WithScale(opts.Scale).
				Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}

func centimetersToInches(cm float64) float64 {
	return cm / 2.54
}
