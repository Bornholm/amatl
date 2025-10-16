package urlx

import (
	"net/url"
	"path/filepath"
)

func Join(dir *url.URL, path string) (*url.URL, error) {
	joined := copy(*dir)

	isOpaque := dir.Opaque != ""

	if isOpaque {
		joined.Path = filepath.Join(dir.Opaque, path)
		joined.Opaque = joined.Path
	} else {
		joined = dir.JoinPath(path)
	}

	return joined, nil
}
