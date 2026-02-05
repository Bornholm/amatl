package file

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/pkg/errors"
)

type Resolver struct {
}

// Resolve implements layout.Resolver.
func (*Resolver) Resolve(ctx context.Context, path resolver.Path) (io.ReadCloser, error) {
	// Get the actual file path, handling file:// URLs
	filePath := path.String()
	scheme := path.Scheme()

	if scheme == "file" {
		// For file:// URLs, we need to handle both absolute and relative paths
		if u, err := path.URL(); err == nil {
			if u.Host != "" && u.Host != "localhost" {
				// Handle file://host/path format (relative paths like file://testdata/test.txt)
				filePath = u.Host + u.Path
			} else {
				// Handle file:///path format (absolute paths)
				filePath = u.Path
				// On Windows, convert /C:/path to C:/path, then let filepath handle separators
				if len(filePath) > 3 && filePath[0] == '/' && len(filePath) > 2 && filePath[2] == ':' {
					filePath = filePath[1:] // Remove leading slash: /C:/path -> C:/path
				}
				// Convert forward slashes to backslashes on Windows
				filePath = filepath.FromSlash(filePath)
			}
		}
	}

	platform := runtime.GOOS
	switch platform {
	case "windows":
		filePath = strings.ReplaceAll(filePath, "/", string(os.PathSeparator))
	default:
		filePath = strings.ReplaceAll(filePath, "\\", string(os.PathSeparator))
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return file, nil
}

func NewResolver() *Resolver {
	return &Resolver{}
}

var _ resolver.Resolver = &Resolver{}
