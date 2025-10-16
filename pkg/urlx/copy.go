package urlx

import "net/url"

func copy(u url.URL) *url.URL {
	return &u
}
