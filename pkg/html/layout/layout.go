package layout

import (
	"context"
	"html/template"
	"io"
	"net/url"

	"github.com/Bornholm/amatl/pkg/html/layout/resolver/amatl"
	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/pkg/errors"
)

type layoutData struct {
	Vars    map[string]any
	Meta    map[string]any
	Body    template.HTML
	Context context.Context
}

func Render(ctx context.Context, w io.Writer, body []byte, funcs ...OptionFunc) error {
	opts := NewLayoutOptions(funcs...)

	url, err := url.Parse(opts.RawURL)
	if err != nil {
		return errors.WithStack(err)
	}

	reader, err := opts.Resolver.Resolve(ctx, url)
	if err != nil {
		return errors.WithStack(err)
	}

	defer func() {
		if err := reader.Close(); err != nil {
			panic(errors.WithStack(err))
		}
	}()

	rawTmpl, err := io.ReadAll(reader)
	if err != nil {
		return errors.WithStack(err)
	}

	layout, err := template.New("").Funcs(opts.Funcs).Parse(string(rawTmpl))
	if err != nil {
		return errors.WithStack(err)
	}

	ctx = resolver.WithWorkDir(ctx, url.JoinPath("../"))

	data := &layoutData{
		Vars:    opts.Vars,
		Meta:    opts.Meta,
		Body:    template.HTML(body),
		Context: ctx,
	}

	if err := layout.Execute(w, data); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

type LayoutOptions struct {
	RawURL   string
	Vars     map[string]any
	Meta     map[string]any
	Resolver resolver.Resolver
	Funcs    template.FuncMap
}

type OptionFunc func(opts *LayoutOptions)

const DefaultRawURL = "amatl://document.html"

func NewLayoutOptions(funcs ...OptionFunc) *LayoutOptions {
	resolver := resolver.DefaultResolver.Extend(
		func() (scheme string, resolver resolver.Resolver) {
			return amatl.Scheme, amatl.NewResolver()
		},
	)
	opts := &LayoutOptions{
		RawURL:   DefaultRawURL,
		Vars:     map[string]any{},
		Meta:     map[string]any{},
		Resolver: resolver,
		Funcs:    DefaultFuncs(resolver),
	}

	for _, fn := range funcs {
		fn(opts)
	}

	return opts
}

func WithVar(key string, value any) OptionFunc {
	return func(opts *LayoutOptions) {
		opts.Vars[key] = value
	}
}

func WithVars(vars map[string]any) OptionFunc {
	return func(opts *LayoutOptions) {
		for key, value := range vars {
			opts.Vars[key] = value
		}
	}
}

func WithMeta(meta map[string]any) OptionFunc {
	return func(opts *LayoutOptions) {
		for key, value := range meta {
			opts.Meta[key] = value
		}
	}
}

func WithURL(rawURL string) OptionFunc {
	return func(opts *LayoutOptions) {
		opts.RawURL = rawURL
	}
}

func WithResolver(resolver resolver.Resolver) OptionFunc {
	return func(opts *LayoutOptions) {
		opts.Resolver = resolver
	}
}
