package layout

import (
	"html/template"
	"net/url"

	"github.com/pkg/errors"
)

type Registry struct {
	resolvers map[string]Resolver
}

// Resolve implements Resolver.
func (r *Registry) Resolve(url *url.URL) (*template.Template, error) {
	resolver, exists := r.resolvers[url.Scheme]
	if !exists {
		return nil, errors.Wrapf(ErrSchemeNotRegistered, "could not resolve scheme '%s'", url.Scheme)
	}

	tmpl, err := resolver.Resolve(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return tmpl, nil
}

func (r *Registry) Register(schema string, resolver Resolver) {
	r.resolvers[schema] = resolver
}

func NewRegistry() *Registry {
	return &Registry{
		resolvers: make(map[string]Resolver),
	}
}

var _ Resolver = &Registry{}
