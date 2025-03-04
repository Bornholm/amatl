package resolver

import (
	"context"
	"io"
	"net/url"

	"github.com/pkg/errors"
)

type Registry struct {
	resolvers map[string]Resolver
}

// Resolve implements Resolver.
func (r *Registry) Resolve(ctx context.Context, url *url.URL) (io.ReadCloser, error) {
	resolver, exists := r.resolvers[url.Scheme]
	if !exists {
		return nil, errors.Wrapf(ErrSchemeNotRegistered, "could not resolve scheme '%s'", url.Scheme)
	}

	reader, err := resolver.Resolve(ctx, url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return reader, nil
}

func (r *Registry) Register(scheme string, resolver Resolver) {
	r.resolvers[scheme] = resolver
}

func (r *Registry) Extend(extensions ...func() (scheme string, resolver Resolver)) *Registry {
	registry := NewRegistry()
	for scheme, resolver := range r.resolvers {
		registry.Register(scheme, resolver)
	}

	for _, ext := range extensions {
		scheme, resolver := ext()
		registry.Register(scheme, resolver)
	}

	return registry
}

func NewRegistry() *Registry {
	return &Registry{
		resolvers: make(map[string]Resolver),
	}
}

var _ Resolver = &Registry{}
