package urlx

import (
	"net/url"
	"path/filepath"

	"github.com/pkg/errors"
)

func Dir(file *url.URL) (*url.URL, error) {
	dir := copy(*file)

	sourcePath := file.Path
	isOpaque := false
	if sourcePath == "" {
		sourcePath = file.Opaque
		isOpaque = true
	}

	dir.Path = filepath.Dir(sourcePath)
	if isOpaque {
		dir.Opaque = dir.Path
	}

	if dir.Scheme == "" {
		dir.Scheme = "file"

		var absPath string
		var err error

		if isOpaque {
			absPath, err = filepath.Abs(dir.Opaque)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		} else {
			absPath, err = filepath.Abs(dir.Path)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		}

		dir.Path = absPath
		if isOpaque {
			dir.Opaque = absPath
		}
	}

	return dir, nil
}
