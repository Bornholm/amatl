package file

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/Bornholm/amatl/pkg/transform"
	"github.com/pkg/errors"
)

func TestResolver(t *testing.T) {
	res := NewResolver()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("%+v", errors.WithStack(err))
	}

	filepath.Join()

	urls := []string{
		"../file/testdata/test.txt",
		"file://testdata/test.txt",
		filepath.Join(workingDirectory, "testdata", "test.txt"),
		// Relative windows path
		"..\\file\\testdata\\test.txt",
	}

	resolve := func(path string) {
		t.Run(path, func(t *testing.T) {
			reader, err := res.Resolve(ctx, resolver.Path(path))
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
