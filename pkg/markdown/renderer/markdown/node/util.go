package node

import (
	"io"
	"strconv"
	"unsafe"

	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/yuin/goldmark/ast"
)

func listItemMarkerChars(tnode *ast.ListItem) []byte {
	parList := tnode.Parent().(*ast.List)
	if parList.IsOrdered() {
		cnt := 1
		if parList.Start != 0 {
			cnt = parList.Start
		}
		s := tnode.PreviousSibling()
		for s != nil {
			cnt++
			s = s.PreviousSibling()
		}
		return append(strconv.AppendInt(nil, int64(cnt), 10), parList.Marker, ' ')
	}
	return []byte{parList.Marker, markdown.SpaceChar[0]}
}

func noAllocString(buf []byte) string {
	return *(*string)(unsafe.Pointer(&buf))
}

// writeClean writes the given byte slice to the writer
// replacing consecutive spaces, newlines, and tabs
// with single spaces.
func writeClean(w io.Writer, bs []byte) error {
	// This works by scanning the byte slice,
	// and writing sub-slices of bs
	// as we see and skip blank sections.

	var (
		// Start of the current sub-slice to be written.
		startIdx int
		// Normalized last character we saw:
		// for whitespace, this is ' ',
		// for everything else, it's left as-is.
		p byte
	)

	for idx, q := range bs {
		if q == '\n' || q == '\r' || q == '\t' {
			q = ' '
		}

		if q == ' ' {
			if p != ' ' {
				// Going from non-blank to blank.
				// Write the current sub-slice and the blank.
				if _, err := w.Write(bs[startIdx:idx]); err != nil {
					return err
				}
				if _, err := w.Write(markdown.SpaceChar); err != nil {
					return err
				}
			}
			startIdx = idx + 1
		} else if p == ' ' {
			// Going from blank to non-blank.
			// Start a new sub-slice.
			startIdx = idx
		}
		p = q
	}

	_, err := w.Write(bs[startIdx:])
	return err
}
