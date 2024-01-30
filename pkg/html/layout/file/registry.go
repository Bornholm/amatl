package file

import "github.com/Bornholm/amatl/pkg/html/layout"

func init() {
	layout.Register("file", NewResolver())
}
