package transform

import (
	"golang.org/x/text/transform"
)

// NewlineTransformer is a transform.Transformer that converts Windows line endings (\r\n) to Unix line endings (\n)
type NewlineTransformer struct{}

// NewNewlineTransformer creates a new NewlineTransformer
func NewNewlineTransformer() *NewlineTransformer {
	return &NewlineTransformer{}
}

// Transform implements the transform.Transformer interface
func (t *NewlineTransformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	for nSrc < len(src) {
		// Check if we have enough space in dst
		if nDst >= len(dst) {
			err = transform.ErrShortDst
			return
		}

		// Look for \r\n sequence
		if nSrc < len(src)-1 && src[nSrc] == '\r' && src[nSrc+1] == '\n' {
			// Replace \r\n with \n
			dst[nDst] = '\n'
			nDst++
			nSrc += 2 // Skip both \r and \n
		} else {
			// Copy byte as-is
			dst[nDst] = src[nSrc]
			nDst++
			nSrc++
		}
	}

	return nDst, nSrc, nil
}

// Reset implements the transform.Transformer interface
func (t *NewlineTransformer) Reset() {
	// No state to reset for this transformer
}

// WindowsToUnixNewlines returns a transform.Transformer that converts Windows line endings to Unix line endings
func WindowsToUnixNewlines() transform.Transformer {
	return NewNewlineTransformer()
}
