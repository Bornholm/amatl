package pipeline

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yuin/goldmark/parser"
)

const contextKey parser.ContextKey = iota

func FromParserContext(pc parser.Context) (context.Context, error) {
	raw := pc.Get(contextKey)
	if raw == nil {
		return nil, errors.New("context not found")
	}

	ctx, ok := raw.(context.Context)
	if !ok {
		return nil, errors.Errorf("unexpected value type '%T'", raw)
	}

	return ctx, nil
}

func WithContext(ctx context.Context, pc parser.Context) parser.Context {
	pc.Set(contextKey, ctx)
	return pc
}
