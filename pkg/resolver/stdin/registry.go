package stdin

import "github.com/Bornholm/amatl/pkg/resolver"

const Scheme = "stdin"

func init() {
	resolver.Register(Scheme, NewResolver())
}
