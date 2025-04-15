package pipeline

import (
	"net/url"
	"path/filepath"

	"github.com/pkg/errors"
)

func SourceDir(sourceURL *url.URL) (*url.URL, error) {
	sourceDir := func(u url.URL) *url.URL {
		return &u
	}(*sourceURL)

	sourceDir.Path = filepath.Dir(sourceURL.Path)

	if sourceDir.Scheme == "" {
		sourceDir.Scheme = "file"

		path, err := filepath.Abs(sourceDir.Path)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		sourceDir.Path = path
	}

	return sourceDir, nil
}
