package directive

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type InlineParser struct {
}

// Parse implements parser.InlineParser.
func (*InlineParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, segment := block.PeekLine()
	stop := findDirectiveEnd(line)
	if stop < 0 {
		return nil
	}

	if stop >= len(line) {
		return nil
	}

	value := ast.NewTextSegment(text.NewSegment(segment.Start+1, segment.Start+stop))

	directive := parseDirective(line[:stop+1], value)
	if directive == nil {
		return nil
	}

	directive.SetParent(parent)

	block.Advance(stop + 1)
	return directive
}

// Trigger implements parser.InlineParser.
func (*InlineParser) Trigger() []byte {
	return []byte(":")
}

var _ parser.InlineParser = &InlineParser{}

func findDirectiveEnd(b []byte) int {
	end := len(b) - 1
	last := b[end]

	if last != ']' && last != '}' {
		return -1
	}

	return end
}
