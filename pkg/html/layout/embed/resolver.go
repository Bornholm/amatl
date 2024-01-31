package embed

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/url"
	"path/filepath"

	"github.com/Bornholm/amatl/pkg/html/layout"
	"github.com/pkg/errors"
)

var (
	//go:embed templates/*.html
	templateFs embed.FS
	templates  *template.Template
)

const templatePattern = "templates/*.html"

func init() {
	tmpl, err := template.New("").ParseFS(templateFs, templatePattern)
	if err != nil {
		panic(errors.Wrap(err, "could not parse embedded layout templates"))
	}

	templates = tmpl
}

func Available() []string {

	filenames, err := fs.Glob(templateFs, templatePattern)
	if err != nil {
		panic(errors.WithStack(err))
	}

	available := make([]string, 0, len(filenames))

	for _, f := range filenames {
		available = append(available, fmt.Sprintf("embed://%s", filepath.Base(f)))
	}

	return available
}

type Resolver struct {
}

// Resolve implements layout.Resolver.
func (*Resolver) Resolve(url *url.URL) (*template.Template, error) {
	name := url.Host

	for _, tmpl := range templates.Templates() {
		if tmpl.Name() != name {
			continue
		}

		return tmpl, nil
	}

	return nil, errors.WithStack(layout.ErrTemplateNotFound)
}

func NewResolver() *Resolver {
	return &Resolver{}
}

var _ layout.Resolver = &Resolver{}
