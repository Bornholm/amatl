package http

import "github.com/Bornholm/amatl/pkg/resolver"

const Scheme = "http"
const SchemeAlt = "https"

func init() {
	httpResolver := NewResolver()
	resolver.Register("http", httpResolver)
	resolver.Register("https", httpResolver)
}
