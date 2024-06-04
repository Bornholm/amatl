package stdin

import (
	"context"
	"io"
	"net/url"
	"os"

	"github.com/Bornholm/amatl/pkg/resolver"
)

type Resolver struct {
}

// Resolve implements layout.Resolver.
func (*Resolver) Resolve(ctx context.Context, url *url.URL) (io.ReadCloser, error) {
	return os.Stdin, nil
}

func NewResolver() *Resolver {
	return &Resolver{}
}

var _ resolver.Resolver = &Resolver{}
