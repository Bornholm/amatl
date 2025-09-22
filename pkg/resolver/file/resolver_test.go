package file

import (
	"context"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/Bornholm/amatl/pkg/transform"
	"github.com/pkg/errors"
)

func TestResolver(t *testing.T) {
	resolver := NewResolver()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("%+v", errors.WithStack(err))
	}

	filepath.Join()

	urls := []string{
		"./testdata/test.txt",
		"file://testdata/test.txt",
		filepath.Join(workingDirectory, "testdata/test.txt"),
	}

	resolve := func(rawURL string) {
		t.Run(rawURL, func(t *testing.T) {
			u, err := url.Parse(rawURL)
			if err != nil {
				t.Errorf("%+v", errors.WithStack(err))
				return
			}

			reader, err := resolver.Resolve(ctx, u)
			if err != nil {
				t.Errorf("%+v", errors.WithStack(err))
				return
			}

			defer reader.Close()

			transformed := transform.NewNewlineReader(reader)

			data, err := io.ReadAll(transformed)
			if err != nil {
				t.Errorf("%+v", errors.WithStack(err))
				return
			}

			if e, g := "foo", string(data); e != g {
				t.Errorf("data: expected '%v', got '%v'", e, g)
			}
		})
	}

	for _, rawURL := range urls {
		resolve(rawURL)
	}
}

func TestWindowsToFilePath(t *testing.T) {
	path := `C:\Users\Docker\Desktop\examples\document\document.md`

	url, err := url.Parse(path)
	if err != nil {
		t.Fatalf("%+v", errors.WithStack(err))
	}

	transformed := toFilePath(url)

	t.Logf("Transformed URL: '%s'", transformed)

	if e, g := `c:\Users\Docker\Desktop\examples\document\document.md`, transformed; e != g {
		t.Errorf("transformed: expected '%s', got '%s'", e, g)
	}
}
