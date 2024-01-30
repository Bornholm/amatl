package layout

import (
	"html/template"
	"net/url"
)

type Resolver interface {
	Resolve(url *url.URL) (*template.Template, error)
}
