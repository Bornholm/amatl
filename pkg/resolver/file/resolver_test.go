package file

import (
	"net/url"
	"testing"

	"github.com/pkg/errors"
)

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
