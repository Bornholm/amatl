package resolver

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type Path string

func (p Path) URL() (*url.URL, error) {
	s := string(p)

	// Quick check: URLs must have a scheme followed by ://
	// This prevents Windows paths like "C:\path" from being parsed as URLs
	if !strings.Contains(s, "://") {
		return nil, errors.New("not a URL")
	}

	u, err := url.ParseRequestURI(s)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Additional validation: scheme must be at least 2 characters
	// This prevents single-letter schemes like "c:" (Windows drives)
	if len(u.Scheme) < 2 {
		return nil, errors.New("invalid URL scheme")
	}

	return u, nil
}

// Scheme returns the scheme of the path if it's a URL, empty string otherwise
func (p Path) Scheme() string {
	if u, err := p.URL(); err == nil {
		return u.Scheme
	}
	return ""
}

// IsURL returns true if the path is a valid URL
func (p Path) IsURL() bool {
	_, err := p.URL()
	return err == nil
}

// WithAuth returns a new Path with basic authentication credentials
func (p Path) WithAuth(username, password string) Path {
	if !p.IsURL() {
		return p
	}

	u, err := p.URL()
	if err != nil {
		return p
	}

	// Clone the URL to avoid modifying the original
	newURL := *u
	if username != "" || password != "" {
		newURL.User = url.UserPassword(username, password)
	}

	return Path(newURL.String())
}

// Host returns the host part of the URL, or empty string if not a URL
func (p Path) Host() string {
	if u, err := p.URL(); err == nil {
		return u.Host
	}
	return ""
}

// URLPath returns the path part of the URL, or the full string if not a URL
func (p Path) URLPath() string {
	if u, err := p.URL(); err == nil {
		return u.Path
	}
	return string(p)
}

func (p Path) Join(paths ...Path) Path {
	strPaths := make([]string, len(paths))
	for i, s := range paths {
		strPaths[i] = s.String()
	}

	if fileURL, err := p.URL(); err == nil {
		return Path(fileURL.JoinPath(strPaths...).String())
	}

	// Create a slice with the base path and all additional paths
	allPaths := append([]string{string(p)}, strPaths...)
	return Path(filepath.Join(allPaths...))
}

func (p Path) IsAbs() bool {
	if p.IsURL() {
		return true
	}

	return filepath.IsAbs(p.String())
}

// Dir returns the directory part of the path
func (p Path) Dir() Path {
	if p.IsURL() {
		u, err := p.URL()
		if err != nil {
			// Fallback to filepath if URL parsing fails
			return Path(filepath.Dir(string(p)))
		}

		// Clone the URL
		dirURL := *u

		// For URLs, always use forward slashes in the path
		urlPath := strings.ReplaceAll(u.Path, "\\", "/")
		dirURL.Path = filepath.ToSlash(filepath.Dir(urlPath))

		// Handle opaque URLs (like Windows paths)
		if u.Opaque != "" {
			if strings.Contains(u.Opaque, "\\") {
				dirURL.Opaque = filepath.Dir(u.Opaque)
				dirURL.Path = strings.ReplaceAll(dirURL.Opaque, "\\", "/")
			} else {
				dirURL.Path = filepath.ToSlash(filepath.Dir(u.Opaque))
				dirURL.Opaque = dirURL.Path
			}
		}

		return Path(dirURL.String())
	}

	return Path(filepath.Dir(string(p)))
}

func (p Path) Base() Path {
	if p.IsURL() {
		u, err := p.URL()
		if err != nil {
			// Fallback to filepath if URL parsing fails
			return Path(filepath.Base(string(p)))
		}

		// Clone the URL
		baseURL := *u

		// For URLs, always use forward slashes in the path
		urlPath := strings.ReplaceAll(u.Path, "\\", "/")
		baseURL.Path = filepath.ToSlash(filepath.Base(urlPath))

		// Handle opaque URLs (like Windows paths)
		if u.Opaque != "" {
			if strings.Contains(u.Opaque, "\\") {
				baseURL.Opaque = filepath.Dir(u.Opaque)
				baseURL.Path = strings.ReplaceAll(baseURL.Opaque, "\\", "/")
			} else {
				baseURL.Path = filepath.ToSlash(filepath.Base(u.Opaque))
				baseURL.Opaque = baseURL.Path
			}
		}

		return Path(baseURL.String())
	}

	return Path(filepath.Base(string(p)))
}

// JoinPath joins the current path with additional path segments
func (p Path) JoinPath(paths ...string) Path {
	if p.IsURL() {
		u, err := p.URL()
		if err != nil {
			// Fallback to filepath if URL parsing fails
			return Path(filepath.Join(append([]string{string(p)}, paths...)...))
		}

		// Use URL's JoinPath method
		newURL := u.JoinPath(paths...)
		return Path(newURL.String())
	}

	// For file paths, use filepath.Join
	return Path(filepath.Join(append([]string{string(p)}, paths...)...))
}

func (p Path) String() string {
	return string(p)
}
