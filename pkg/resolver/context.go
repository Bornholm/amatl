package resolver

import (
	"context"
	"net/url"
)

type contextKey string

const (
	contextKeyWorkDir  contextKey = "workdir"
	contextKeyResolver contextKey = "resolver"
)

func WithWorkDir(ctx context.Context, url *url.URL) context.Context {
	return context.WithValue(ctx, contextKeyWorkDir, url)
}

func ContextWorkDir(ctx context.Context) *url.URL {
	workDir, ok := ctx.Value(contextKeyWorkDir).(*url.URL)
	if !ok {
		return nil
	}

	return workDir
}

func WithResolver(ctx context.Context, resolver Resolver) context.Context {
	return context.WithValue(ctx, contextKeyResolver, resolver)
}

func ContextResolver(ctx context.Context) Resolver {
	resolver, ok := ctx.Value(contextKeyResolver).(Resolver)
	if !ok {
		return nil
	}

	return resolver
}
