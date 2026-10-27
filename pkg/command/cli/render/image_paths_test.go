package render

import (
	"regexp"
	"testing"
)

// TestIncludeImagePaths checks that images of an included file are found
// whatever the way their path is written: between angle brackets, percent
// encoded, or with escaped parentheses
func TestIncludeImagePaths(t *testing.T) {
	html := renderContainerTestData(t, "testdata/include-image-paths")

	for _, alt := range []string{"angle", "encoded", "parens", "escaped"} {
		img := regexp.MustCompile(`<img src="data:image/svg\+xml;base64,[^"]+" alt="` + alt + `"`)
		if !img.MatchString(html) {
			t.Errorf("expected an embedded '%s' image in html:\n%s", alt, html)
		}
	}
}
