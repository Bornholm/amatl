package amatl

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"path/filepath"

	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/pkg/errors"
)

const Scheme = "amatl"

var (
	//go:embed templates/*.html
	templateFs embed.FS
)

const templatePattern = "templates/*.html"

func Available() []string {
	filenames, err := fs.Glob(templateFs, templatePattern)
	if err != nil {
		panic(errors.WithStack(err))
	}

	available := make([]string, 0, len(filenames))

	for _, f := range filenames {
		available = append(available, fmt.Sprintf("amatl://%s", filepath.Base(f)))
	}

	return available
}

type Resolver struct {
}

// Resolve implements layout.Resolver.
func (*Resolver) Resolve(ctx context.Context, url *url.URL) (io.ReadCloser, error) {
	filename := url.Host

	file, err := templateFs.Open(filepath.Join("templates", filename))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return file, nil
}

func NewResolver() *Resolver {
	return &Resolver{}
}

var _ resolver.Resolver = &Resolver{}
