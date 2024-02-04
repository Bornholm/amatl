package base

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/url"
	"path/filepath"

	"github.com/Bornholm/amatl/pkg/html/layout"
	"github.com/pkg/errors"
)

var (
	//go:embed templates/*.html
	templateFs embed.FS
)

const templatePattern = "templates/*.html"

func Available() []string {
	filenames, err := fs.Glob(templateFs, templatePattern)
	if err != nil {
		panic(errors.WithStack(err))
	}

	available := make([]string, 0, len(filenames))

	for _, f := range filenames {
		available = append(available, fmt.Sprintf("base://%s", filepath.Base(f)))
	}

	return available
}

type Resolver struct {
}

// Resolve implements layout.Resolver.
func (*Resolver) Resolve(url *url.URL, funcs template.FuncMap) (*template.Template, error) {
	templates, err := template.New("").Funcs(funcs).ParseFS(templateFs, templatePattern)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse base templates")
	}

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