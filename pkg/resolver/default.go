package resolver

import (
	"context"
	"io"
	"net/url"

	"github.com/pkg/errors"
)

var DefaultResolver = NewRegistry()

func Register(scheme string, resolver Resolver) {
	DefaultResolver.Register(scheme, resolver)
}

func Resolve(ctx context.Context, url *url.URL) (io.ReadCloser, error) {
	reader, err := DefaultResolver.Resolve(ctx, url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return reader, nil
}
