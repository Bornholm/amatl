package http

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/pkg/errors"
)

type Resolver struct {
}

// Resolve implements layout.Resolver.
func (*Resolver) Resolve(ctx context.Context, url *url.URL) (io.ReadCloser, error) {
	res, err := http.Get(url.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res.Body, nil
}

func NewResolver() *Resolver {
	return &Resolver{}
}

var _ resolver.Resolver = &Resolver{}
