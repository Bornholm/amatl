package file

import (
	"context"
	"io"
	"net/url"
	"os"

	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/pkg/errors"
)

type Resolver struct {
}

// Resolve implements layout.Resolver.
func (*Resolver) Resolve(ctx context.Context, url *url.URL) (io.ReadCloser, error) {
	path := url.Host + url.Path

	file, err := os.Open(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return file, nil
}

func NewResolver() *Resolver {
	return &Resolver{}
}

var _ resolver.Resolver = &Resolver{}
