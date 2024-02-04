package layout

import (
	"html/template"
	"net/url"

	"github.com/pkg/errors"
)

var (
	defaultRegistry = NewRegistry()
)

func Register(scheme string, resolver Resolver) {
	defaultRegistry.Register(scheme, resolver)
}

func Resolve(url *url.URL) (*template.Template, error) {
	tmpl, err := defaultRegistry.Resolve(url, DefaultFuncs())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return tmpl, nil
}
