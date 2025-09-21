package transform

import (
	"io"

	"golang.org/x/text/transform"
)

func NewNewlineReader(r io.Reader) io.Reader {
	return transform.NewReader(r, NewNewlineTransformer())
}
