package resolver

import (
	"context"
	"io"
	"net/url"
)

type Resolver interface {
	Resolve(ctx context.Context, url *url.URL) (io.ReadCloser, error)
}
