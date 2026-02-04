package resolver

import (
	"context"
	"io"
	"log/slog"

	"github.com/pkg/errors"
)

type Registry struct {
	resolvers       map[string]Resolver
	defaultResolver string
}

// Resolve implements Resolver.
func (r *Registry) Resolve(ctx context.Context, path Path) (io.ReadCloser, error) {
	ctx = WithResolver(ctx, r)

	// Handle relative paths with working directory first
	workDir := ContextWorkDir(ctx)
	resolvedPath := path
	if workDir != "" && !path.IsAbs() {
		resolvedPath = workDir.JoinPath(path.String())
		slog.DebugContext(ctx, "using workdir", slog.String("original_path", path.String()), slog.String("workdir", workDir.String()), slog.String("joined_path", resolvedPath.String()))

	}

	// Now determine the scheme from the resolved path
	scheme := resolvedPath.Scheme()

	resolver, exists := r.resolvers[scheme]
	if !exists {
		if r.defaultResolver != "" {
			resolver = r.resolvers[r.defaultResolver]
		}

		if resolver == nil {
			return nil, errors.Wrapf(ErrSchemeNotRegistered, "could not resolve path '%s'", resolvedPath)
		}
	}

	slog.DebugContext(ctx, "resolving path", slog.String("path", resolvedPath.String()))

	reader, err := resolver.Resolve(ctx, resolvedPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return reader, nil
}

func (r *Registry) Register(scheme string, resolver Resolver) {
	r.resolvers[scheme] = resolver
}

func (r *Registry) SetDefault(scheme string) {
	r.defaultResolver = scheme
}

func (r *Registry) Extend(extensions ...func() (scheme string, resolver Resolver)) *Registry {
	registry := NewRegistry()
	for scheme, resolver := range r.resolvers {
		registry.Register(scheme, resolver)
	}

	registry.SetDefault(r.defaultResolver)

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
