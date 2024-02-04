package base

import "github.com/Bornholm/amatl/pkg/html/layout"

func init() {
	layout.Register("base", NewResolver())
}
