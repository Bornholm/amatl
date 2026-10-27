package render

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// TestIncludeRemoteImagePaths checks that images of a file included over
// HTTP are fetched with their path encoded once
func TestIncludeRemoteImagePaths(t *testing.T) {
	server := httptest.NewServer(http.FileServer(http.Dir("testdata/include-remote")))
	defer server.Close()

	dir := t.TempDir()

	layout, err := os.ReadFile("testdata/container/layout.html")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	files := map[string]string{
		"main.md":     `:include{url="` + server.URL + `/sub/part.md"}` + "\n",
		"layout.html": string(layout),
	}

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatalf("%+v", err)
		}
	}

	html := renderContainerTestData(t, dir)

	for _, alt := range []string{"space", "accent"} {
		img := regexp.MustCompile(`<img src="data:image/svg\+xml;base64,[^"]+" alt="` + alt + `"`)
		if !img.MatchString(html) {
			t.Errorf("expected an embedded '%s' image in html:\n%s", alt, html)
		}
	}
}
