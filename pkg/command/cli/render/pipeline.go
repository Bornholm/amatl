package render

import (
	"bytes"
	"context"
	"log/slog"
	"net/url"
	"sync"
	"text/template"
	"time"

	"github.com/Bornholm/amatl/pkg/html/layout"
	"github.com/Bornholm/amatl/pkg/pipeline"
	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/Masterminds/sprig/v3"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

const (
	attrMeta = "meta"
)

type TemplateTransformerOptions struct {
	Vars  map[string]any
	Funcs template.FuncMap
}

type TemplateTransformerOptionFunc func(opts *TemplateTransformerOptions)

func NewTemplateTransformerOptions(funcs ...TemplateTransformerOptionFunc) *TemplateTransformerOptions {
	opts := &TemplateTransformerOptions{
		Vars:  map[string]any{},
		Funcs: sprig.FuncMap(),
	}
	for _, fn := range funcs {
		fn(opts)
	}
	return opts
}

func WithVars(vars map[string]any) TemplateTransformerOptionFunc {
	return func(opts *TemplateTransformerOptions) {
		for key, value := range vars {
			opts.Vars[key] = value
		}
	}
}

func WithFuncs(funcs template.FuncMap) TemplateTransformerOptionFunc {
	return func(opts *TemplateTransformerOptions) {
		opts.Funcs = funcs
	}
}

func TemplateMiddleware(funcs ...TemplateTransformerOptionFunc) pipeline.Middleware {
	opts := NewTemplateTransformerOptions(funcs...)
	return func(next pipeline.Transformer) pipeline.Transformer {
		return pipeline.TransformerFunc(func(ctx context.Context, payload *pipeline.Payload) error {
			data := payload.GetData()

			tmpl, err := template.New("").Funcs(opts.Funcs).Parse(string(data))
			if err != nil {
				return errors.WithStack(err)
			}

			meta, ok := pipeline.GetAttribute[map[string]any](payload, attrMeta)
			if !ok {
				meta = make(map[string]any)
			}

			vars := struct {
				Vars map[string]any
				Meta map[string]any
			}{
				Vars: opts.Vars,
				Meta: meta,
			}

			var doc bytes.Buffer

			if err := tmpl.Execute(&doc, vars); err != nil {
				return errors.WithStack(err)
			}

			payload.SetData(doc.Bytes())

			if err := next.Transform(ctx, payload); err != nil {
				return errors.WithStack(err)
			}

			return nil
		})
	}
}

type MarkdownTransformerOptions struct {
	SourceURL        *url.URL
	LinkReplacements map[string]string
}

type MarkdownTransformerOptionFunc func(opts *MarkdownTransformerOptions)

func NewMarkdownTransformerOptions(funcs ...MarkdownTransformerOptionFunc) *MarkdownTransformerOptions {
	opts := &MarkdownTransformerOptions{
		LinkReplacements: make(map[string]string),
	}
	for _, fn := range funcs {
		fn(opts)
	}
	return opts
}

func WithSourceURL(sourceURL *url.URL) MarkdownTransformerOptionFunc {
	return func(opts *MarkdownTransformerOptions) {
		opts.SourceURL = sourceURL
	}
}

func WithLinkReplacements(replacements map[string]string) MarkdownTransformerOptionFunc {
	return func(opts *MarkdownTransformerOptions) {
		opts.LinkReplacements = replacements
	}
}

func MarkdownMiddleware(funcs ...MarkdownTransformerOptionFunc) pipeline.Middleware {
	opts := NewMarkdownTransformerOptions(funcs...)
	return func(next pipeline.Transformer) pipeline.Transformer {
		return pipeline.TransformerFunc(func(ctx context.Context, payload *pipeline.Payload) error {
			slog.DebugContext(ctx, "entering markdown middleware")

			data := payload.GetData()

			reader := text.NewReader(data)

			sourceDir, err := pipeline.SourceDir(opts.SourceURL)
			if err != nil {
				return errors.WithStack(err)
			}

			ctx = resolver.WithWorkDir(ctx, sourceDir)

			parse := newParser(opts.SourceURL, ParserOptions{
				EmbedLinkedResources: false,
				LinkReplacements:     opts.LinkReplacements,
			})
			render := newMarkdownRenderer()

			slog.DebugContext(ctx, "parsing markdown file")

			pc := parser.NewContext()
			pc = pipeline.WithContext(ctx, pc)

			document := parse.Parse(reader, parser.WithContext(pc))

			var doc bytes.Buffer

			if err := render.Render(&doc, data, document); err != nil {
				return errors.WithStack(err)
			}

			payload.SetData(doc.Bytes())

			meta := document.OwnerDocument().Meta()
			payload.SetAttribute(attrMeta, meta)

			if err := next.Transform(ctx, payload); err != nil {
				return errors.WithStack(err)
			}

			slog.DebugContext(ctx, "exiting markdown middleware")

			return nil
		})
	}
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
		LayoutURL:                  layout.DefaultRawURL,
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

func HTMLMiddleware(funcs ...HTMLTransformerOptionFunc) pipeline.Middleware {
	opts := NewHTMLTransformerOptions(funcs...)
	return func(next pipeline.Transformer) pipeline.Transformer {
		return pipeline.TransformerFunc(func(ctx context.Context, payload *pipeline.Payload) error {
			slog.DebugContext(ctx, "entering html middleware")

			data := payload.GetData()

			reader := text.NewReader(data)

			sourceDir, err := pipeline.SourceDir(opts.SourceURL)
			if err != nil {
				return errors.WithStack(err)
			}

			ctx = resolver.WithWorkDir(ctx, sourceDir)

			parse := newParser(opts.SourceURL, ParserOptions{
				EmbedLinkedResources: true,
				LinkReplacements:     opts.LinkReplacements,
			})
			pc := parser.NewContext()
			pc = pipeline.WithContext(ctx, pc)

			slog.DebugContext(ctx, "parsing markdown file")

			document := parse.Parse(reader, parser.WithContext(pc))

			meta, ok := pipeline.GetAttribute[map[string]any](payload, attrMeta)
			if !ok {
				meta = make(map[string]any)
			}

			render := newHTMLRenderer()

			var body bytes.Buffer

			if err := render.Render(&body, data, document); err != nil {
				return errors.WithStack(err)
			}

			var doc bytes.Buffer

			layoutURL, err := url.Parse(opts.LayoutURL)
			if err != nil {
				return errors.WithStack(err)
			}

			workDir, err := pipeline.SourceDir(layoutURL)
			if err != nil {
				return errors.WithStack(err)
			}

			ctx = resolver.WithWorkDir(ctx, workDir)

			slog.DebugContext(ctx, "rendering html layout", slog.String("layout", opts.LayoutURL))

			err = layout.Render(
				ctx, &doc, body.Bytes(),
				layout.WithURL(opts.LayoutURL),
				layout.WithVars(opts.LayoutVars),
				layout.WithMeta(meta),
			)
			if err != nil {
				return errors.WithStack(err)
			}

			payload.SetData(doc.Bytes())

			if err := next.Transform(ctx, payload); err != nil {
				return errors.WithStack(err)
			}

			slog.DebugContext(ctx, "exiting html middleware")

			return nil
		})
	}
}

type PDFTransformerOptions struct {
	MarginTop    float64
	MarginLeft   float64
	MarginRight  float64
	MarginBottom float64
	Scale        float64
	Background   bool
	Timeout      time.Duration
	ExecPath     string
}

const (
	DefaultPDFMargin     float64       = 1
	DefaultPDFScale      float64       = 1
	DefaultPDFTimeout    time.Duration = time.Minute
	DefaultPDFBackground bool          = true
	DefaultPDFExecPath   string        = ""
)

type PDFTransformerOptionFunc func(opts *PDFTransformerOptions)

func NewPDFTransformerOptions(funcs ...PDFTransformerOptionFunc) *PDFTransformerOptions {
	opts := &PDFTransformerOptions{
		MarginTop:    DefaultPDFMargin,
		MarginLeft:   DefaultPDFMargin,
		MarginRight:  DefaultPDFMargin,
		MarginBottom: DefaultPDFMargin,
		Scale:        DefaultPDFScale,
		Background:   DefaultPDFBackground,
		Timeout:      DefaultPDFTimeout,
		ExecPath:     DefaultPDFExecPath,
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

func WithTimeout(timeout time.Duration) PDFTransformerOptionFunc {
	return func(opts *PDFTransformerOptions) {
		opts.Timeout = timeout
	}
}

func WithBackground(background bool) PDFTransformerOptionFunc {
	return func(opts *PDFTransformerOptions) {
		opts.Background = background
	}
}

func WithExecPath(execPath string) PDFTransformerOptionFunc {
	return func(opts *PDFTransformerOptions) {
		opts.ExecPath = execPath
	}
}

func PDFMiddleware(funcs ...PDFTransformerOptionFunc) pipeline.Middleware {
	opts := NewPDFTransformerOptions(funcs...)

	return func(next pipeline.Transformer) pipeline.Transformer {
		return pipeline.TransformerFunc(func(ctx context.Context, payload *pipeline.Payload) error {
			slog.DebugContext(ctx, "entering pdf middleware")

			data := payload.GetData()

			allocatorOptions := chromedp.DefaultExecAllocatorOptions[:]

			if opts.ExecPath != "" {
				allocatorOptions = append(allocatorOptions, chromedp.ExecPath(opts.ExecPath))
			}

			allocatorCtx, allocatorCtxCancel := chromedp.NewExecAllocator(ctx,
				allocatorOptions...,
			)
			defer allocatorCtxCancel()

			timeoutCtx, timeoutCancel := context.WithTimeout(allocatorCtx, opts.Timeout)
			defer timeoutCancel()

			ctx, cancel := chromedp.NewContext(timeoutCtx)
			defer cancel()

			var output []byte

			slog.DebugContext(ctx, "rendering pdf with chrome", slog.Duration("timeout", opts.Timeout))

			if err := chromedp.Run(ctx, printToPDF(data, &output, opts)); err != nil {
				return errors.Wrap(err, "could not execute chrome")
			}

			payload.SetData(output)

			if err := next.Transform(ctx, payload); err != nil {
				return errors.WithStack(err)
			}

			slog.DebugContext(ctx, "exiting pdf middleware")

			return nil
		})
	}
}

func printToPDF(html []byte, res *[]byte, opts *PDFTransformerOptions) chromedp.Tasks {
	return chromedp.Tasks{
		enableLifeCycleEvents(),
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			lctx, cancel := context.WithCancel(ctx)
			defer cancel()
			var wg sync.WaitGroup
			wg.Add(2)
			chromedp.ListenTarget(lctx, func(e any) {
				switch evt := e.(type) {
				case *page.EventLoadEventFired:
					wg.Done()
				case *page.EventLifecycleEvent:
					if evt.Name == "networkIdle" {
						wg.Done()
					}
				}
			})
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return errors.WithStack(err)
			}
			if err := page.SetDocumentContent(frameTree.Frame.ID, string(html)).Do(ctx); err != nil {
				return errors.WithStack(err)
			}
			wg.Wait()
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithDisplayHeaderFooter(false).
				WithPreferCSSPageSize(true).
				WithMarginRight(centimetersToInches(opts.MarginRight)).
				WithMarginTop(centimetersToInches(opts.MarginTop)).
				WithMarginBottom(centimetersToInches(opts.MarginBottom)).
				WithMarginLeft(centimetersToInches(opts.MarginLeft)).
				WithPrintBackground(opts.Background).
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

func enableLifeCycleEvents() chromedp.ActionFunc {
	return func(ctx context.Context) error {
		err := page.Enable().Do(ctx)
		if err != nil {
			return err
		}
		err = page.SetLifecycleEventsEnabled(true).Do(ctx)
		if err != nil {
			return err
		}
		return nil
	}
}

func centimetersToInches(cm float64) float64 {
	return cm / 2.54
}

func ToggleableMiddleware(t pipeline.Transformer, enabled bool) pipeline.Middleware {
	return func(next pipeline.Transformer) pipeline.Transformer {
		return pipeline.TransformerFunc(func(ctx context.Context, payload *pipeline.Payload) error {
			if enabled {
				if err := t.Transform(ctx, payload); err != nil {
					return errors.WithStack(err)
				}
			}

			if err := next.Transform(ctx, payload); err != nil {
				return errors.WithStack(err)
			}

			return nil
		})
	}
}
