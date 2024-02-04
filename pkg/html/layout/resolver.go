package layout

import (
	"html/template"
	"net/url"
)

type Resolver interface {
	Resolve(url *url.URL, funcs template.FuncMap) (*template.Template, error)
}
