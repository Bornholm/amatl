package urlx

import (
	"net/url"
	"path/filepath"
	"strings"

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

	// Fix Windows path handling: convert backslashes to forward slashes
	// before calling filepath.Dir to ensure proper directory extraction
	if isOpaque && strings.Contains(sourcePath, "\\") {
		// This is likely a Windows path in opaque format
		// Convert backslashes to forward slashes for proper directory calculation
		normalizedPath := strings.ReplaceAll(sourcePath, "\\", "/")
		dir.Path = filepath.Dir(normalizedPath)
		// Convert back to backslashes for Windows compatibility
		dir.Path = strings.ReplaceAll(dir.Path, "/", "\\")
	} else {
		dir.Path = filepath.Dir(sourcePath)
	}

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
