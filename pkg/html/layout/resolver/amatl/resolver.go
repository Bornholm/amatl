package amatl

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/pkg/errors"
)

const Scheme = "amatl"

var (
	//go:embed templates/*.html assets/*
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
func (*Resolver) Resolve(ctx context.Context, path resolver.Path) (io.ReadCloser, error) {
	filename := path.Host()

	for _, dir := range []string{"templates", "assets"} {
		file, err := templateFs.Open(dir + "/" + filename)
		if err == nil {
			return file, nil
		}
	}

	return nil, errors.Errorf("amatl: resource %q not found", filename)
}

func NewResolver() *Resolver {
	return &Resolver{}
}

var _ resolver.Resolver = &Resolver{}
