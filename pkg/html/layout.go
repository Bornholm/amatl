package html

import (
	"html/template"
	"io"
	"net/url"

	"github.com/Bornholm/amatl/pkg/html/layout"
	"github.com/pkg/errors"

	// Register layout resolvers
	_ "github.com/Bornholm/amatl/pkg/html/layout/embed"
	_ "github.com/Bornholm/amatl/pkg/html/layout/file"
)

type layoutData struct {
	Vars map[string]any
	Body template.HTML
}

func Layout(w io.Writer, body []byte, funcs ...LayoutOptionFunc) error {
	opts := NewLayoutOptions(funcs...)

	layoutURL, err := url.Parse(opts.RawLayoutURL)
	if err != nil {
		return errors.WithStack(err)
	}

	layoutTmpl, err := layout.Resolve(layoutURL)
	if err != nil {
		return errors.WithStack(err)
	}

	data := &layoutData{
		Vars: opts.Vars,
		Body: template.HTML(body),
	}

	if err := layoutTmpl.Execute(w, data); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

type LayoutOptions struct {
	RawLayoutURL string
	Vars         map[string]any
}

type LayoutOptionFunc func(opts *LayoutOptions)

func NewLayoutOptions(funcs ...LayoutOptionFunc) *LayoutOptions {
	opts := &LayoutOptions{
		RawLayoutURL: "embed://document.html",
		Vars:         map[string]any{},
	}

	for _, fn := range funcs {
		fn(opts)
	}

	return opts
}

func WithVar(key string, value any) LayoutOptionFunc {
	return func(opts *LayoutOptions) {
		opts.Vars[key] = value
	}
}

func WithLayout(rawURL string) LayoutOptionFunc {
	return func(opts *LayoutOptions) {
		opts.RawLayoutURL = rawURL
	}
}
