package resolver

import (
	"context"
	"log/slog"
)

type contextKey string

const (
	contextKeyWorkDir  contextKey = "workdir"
	contextKeyResolver contextKey = "resolver"
)

func WithWorkDir(ctx context.Context, path Path) context.Context {
	slog.DebugContext(ctx, "using work dir", slog.String("workdir", path.String()))
	return context.WithValue(ctx, contextKeyWorkDir, path)
}

func ContextWorkDir(ctx context.Context) Path {
	workDir, ok := ctx.Value(contextKeyWorkDir).(Path)
	if !ok {
		return ""
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
