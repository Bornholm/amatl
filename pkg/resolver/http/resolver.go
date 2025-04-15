package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/pkg/errors"
)

type Resolver struct {
}

// Resolve implements layout.Resolver.
func (*Resolver) Resolve(ctx context.Context, u *url.URL) (io.ReadCloser, error) {
	username := os.Getenv("AMATL_HTTP_BASIC_AUTH_USERNAME")
	password := os.Getenv("AMATL_HTTP_BASIC_AUTH_PASSWORD")

	if username != "" || password != "" {
		u = cloneURL(*u)
		u.User = url.UserPassword(username, password)
	}

	var buff bytes.Buffer

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), &buff)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.URL = u

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

func cloneURL(u url.URL) *url.URL {
	return &u
}
