package file

import "github.com/Bornholm/amatl/pkg/resolver"

const Scheme = "file"
const SchemeAlt = ""

func init() {
	resolver.Register(Scheme, NewResolver())
	resolver.Register(SchemeAlt, NewResolver())
}
