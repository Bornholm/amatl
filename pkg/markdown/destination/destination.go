// Package destination converts between the destination of a markdown link
// or image, as written in the source, and the path or URL it points to.
//
// goldmark keeps destinations as written: backslash escapes and percent
// encoding stay in ast.Link.Destination and ast.Image.Destination, and only
// its HTML renderer decodes them.
package destination

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/yuin/goldmark/util"
)

// Decode returns the path or URL the given destination points to.
// Backslash escapes are removed, and the percent encoding of a local path is
// decoded, so ./my%20image.png points to the "my image.png" file. A URL with a
// scheme is returned with its percent encoding, which its resolver handles.
func Decode(raw string) string {
	decoded := string(util.UnescapePunctuations([]byte(raw)))

	if strings.Contains(decoded, "://") {
		return decoded
	}

	if unescaped, err := url.PathUnescape(decoded); err == nil {
		return unescaped
	}

	return decoded
}

// encoder percent-encodes the characters that cannot be written as is in a
// destination, or would be read differently by Decode
var encoder = strings.NewReplacer(
	"%", "%25",
	" ", "%20",
	"\t", "%09",
	"\n", "%0A",
	"\r", "%0D",
	"(", "%28",
	")", "%29",
	"<", "%3C",
	">", "%3E",
	"\\", "%5C",
)

// Encode returns a destination pointing to the given local path, with slash
// separators, that Decode turns back into the same path.
func Encode(path string) string {
	return encoder.Replace(filepath.ToSlash(path))
}

// Join returns a destination pointing to the given decoded relative path,
// resolved against dir. When dir is a URL, the result of its JoinPath is
// already percent-encoded and is returned as is.
func Join(dir resolver.Path, path string) string {
	joined := dir.JoinPath(path)

	if dir.IsURL() {
		return joined.String()
	}

	return Encode(joined.String())
}

// EncodeURLPath percent-encodes the given decoded URL path, '#', '?' and
// non-ASCII characters included, as url.URL.JoinPath does
func EncodeURLPath(path string) string {
	return (&url.URL{Path: path}).EscapedPath()
}

// Location returns the given destination with its escapes and its percent
// encoding removed, URL or not, to compare it with a location written by
// hand
func Location(raw string) string {
	decoded := string(util.UnescapePunctuations([]byte(raw)))

	if unescaped, err := url.PathUnescape(decoded); err == nil {
		return unescaped
	}

	return decoded
}
