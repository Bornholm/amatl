package file

import (
	"html/template"
	"net/url"

	"github.com/Bornholm/amatl/pkg/html/layout"
	"github.com/pkg/errors"
)

type Resolver struct {
}

// Resolve implements layout.Resolver.
func (*Resolver) Resolve(url *url.URL, funcs template.FuncMap) (*template.Template, error) {
	path := url.Host + url.Path

	tmpl, err := template.New("").Funcs(funcs).ParseFiles(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return tmpl, nil
}

func NewResolver() *Resolver {
	return &Resolver{}
}

var _ layout.Resolver = &Resolver{}
