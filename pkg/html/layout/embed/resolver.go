package embed

import (
	"embed"
	"html/template"
	"net/url"

	"github.com/Bornholm/amatl/pkg/html/layout"
	"github.com/pkg/errors"
)

var (
	//go:embed templates/*.html
	fs        embed.FS
	templates *template.Template
)

func init() {
	tmpl, err := template.New("").ParseFS(fs, "templates/*.html")
	if err != nil {
		panic(errors.Wrap(err, "could not parse embedded layout templates"))
	}

	templates = tmpl
}

type Resolver struct {
}

// Resolve implements layout.Resolver.
func (*Resolver) Resolve(url *url.URL) (*template.Template, error) {
	name := url.Host

	for _, tmpl := range templates.Templates() {
		if tmpl.Name() != name {
			continue
		}

		return tmpl, nil
	}

	return nil, errors.WithStack(layout.ErrTemplateNotFound)
}

func NewResolver() *Resolver {
	return &Resolver{}
}

var _ layout.Resolver = &Resolver{}
