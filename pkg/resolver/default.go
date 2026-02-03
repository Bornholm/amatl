package resolver

import (
	"context"
	"io"

	"github.com/pkg/errors"
)

var DefaultResolver = NewRegistry()

func Register(scheme string, resolver Resolver) {
	DefaultResolver.Register(scheme, resolver)
}

func SetDefault(scheme string) {
	DefaultResolver.SetDefault(scheme)
}

func Resolve(ctx context.Context, path string) (io.ReadCloser, error) {
	reader, err := DefaultResolver.Resolve(ctx, Path(path))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return reader, nil
}
