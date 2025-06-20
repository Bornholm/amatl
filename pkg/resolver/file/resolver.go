package file

import (
	"context"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/pkg/errors"
)

type Resolver struct {
}

// Resolve implements layout.Resolver.
func (*Resolver) Resolve(ctx context.Context, url *url.URL) (io.ReadCloser, error) {
	path := toFilePath(url)

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

		path = toFilePath(absURL)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return file, nil
}

func toFilePath(u *url.URL) string {
	var path string

	if u.Opaque != "" {
		path = u.Scheme + `:\` + strings.ReplaceAll(u.Opaque, `\\`, `\`) // Windows path
	} else {
		path = u.Host + u.Path
	}

	return path
}

func NewResolver() *Resolver {
	return &Resolver{}
}

var _ resolver.Resolver = &Resolver{}
