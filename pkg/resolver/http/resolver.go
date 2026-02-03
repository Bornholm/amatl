package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"

	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/pkg/errors"
)

type Resolver struct {
}

// Resolve implements layout.Resolver.
func (*Resolver) Resolve(ctx context.Context, path resolver.Path) (io.ReadCloser, error) {
	scheme := path.Scheme()

	// Only handle HTTP/HTTPS schemes
	if scheme != "http" && scheme != "https" {
		return nil, errors.Errorf("http resolver can only handle http/https schemes, got: %s", scheme)
	}

	username := os.Getenv("AMATL_HTTP_BASIC_AUTH_USERNAME")
	password := os.Getenv("AMATL_HTTP_BASIC_AUTH_PASSWORD")

	// Use Path's WithAuth method to handle authentication
	requestPath := path
	if username != "" || password != "" {
		requestPath = path.WithAuth(username, password)
	}

	var buff bytes.Buffer

	req, err := http.NewRequestWithContext(ctx, "GET", requestPath.String(), &buff)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return nil, errors.Errorf("unexpected http response status '%s' (%d)", res.Status, res.StatusCode)
	}

	return res.Body, nil
}

func NewResolver() *Resolver {
	return &Resolver{}
}

var _ resolver.Resolver = &Resolver{}
