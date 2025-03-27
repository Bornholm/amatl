package file

import (
	"context"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/pkg/errors"
)

type Resolver struct {
}

// Resolve implements layout.Resolver.
func (*Resolver) Resolve(ctx context.Context, url *url.URL) (io.ReadCloser, error) {
	path := url.Host + url.Path

	// Handle relative "urls"
	workDir := resolver.ContextWorkDir(ctx)
	if workDir != nil && !filepath.IsAbs(path) {
		absURL := workDir.JoinPath(path)
		if absURL.Scheme != Scheme && absURL.Scheme != SchemeAlt {
			reader, err := resolver.ContextResolver(ctx).Resolve(ctx, absURL)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			return reader, nil
		}

		path = absURL.Host + absURL.Path
	}

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
