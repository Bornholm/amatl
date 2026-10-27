package node

import (
	"bytes"
	"io"
	"strconv"
	"unsafe"

	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
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

// formatDestination returns the destination of a link or an image as it must
// be written between parentheses. A destination holding a space, a '<' or
// unbalanced parentheses is only valid between angle brackets, where an
// unescaped '<' or '>' must be escaped. Line breaks are valid in neither
// form and are percent-encoded.
func formatDestination(destination []byte) []byte {
	destination = bytes.ReplaceAll(destination, []byte("\n"), []byte("%0A"))
	destination = bytes.ReplaceAll(destination, []byte("\r"), []byte("%0D"))

	if !needsAngleBrackets(destination) {
		return destination
	}

	formatted := make([]byte, 0, len(destination)+4)
	formatted = append(formatted, '<')

	for i := 0; i < len(destination); i++ {
		c := destination[i]

		if isEscape(destination, i) {
			formatted = append(formatted, c, destination[i+1])
			i++
			continue
		}

		if c == '<' || c == '>' {
			formatted = append(formatted, '\\')
		}

		formatted = append(formatted, c)
	}

	return append(formatted, '>')
}

func needsAngleBrackets(destination []byte) bool {
	depth := 0

	for i := 0; i < len(destination); i++ {
		if isEscape(destination, i) {
			i++
			continue
		}

		switch c := destination[i]; {
		case c <= ' ' || c == '<':
			return true
		case c == '(':
			depth++
		case c == ')':
			if depth--; depth < 0 {
				return true
			}
		}
	}

	return depth != 0
}

// isEscape reports whether the byte at i is a backslash escaping the next
// one: CommonMark only treats a backslash before an ASCII punctuation
// character as an escape
func isEscape(destination []byte, i int) bool {
	return destination[i] == '\\' && i+1 < len(destination) && util.IsPunct(destination[i+1])
}
