package resolver

import (
	"context"
	"io"
)

type Resolver interface {
	Resolve(ctx context.Context, path Path) (io.ReadCloser, error)
}
