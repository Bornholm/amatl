package embed

import "github.com/Bornholm/amatl/pkg/html/layout"

func init() {
	layout.Register("embed", NewResolver())
}
