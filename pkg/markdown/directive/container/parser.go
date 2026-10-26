package container

import (
	"github.com/Bornholm/amatl/pkg/markdown/directive"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

const (
	fenceChar      = ':'
	minFenceLength = 3
)

// BlockParser parses container directives.
//
// A container is opened by a fence of at least three colons followed by
// a name and optional attributes, and closed by a fence of at least the same
// length. A closing fence closes the innermost opened container, so containers
// can be nested with fences of the same length, unless the fence is less
// indented than that container's opening fence (e.g. a fence at the beginning
// of the line when the innermost container is inside a list item): it then
// closes the outer container and, implicitly, everything opened inside it.
type BlockParser struct {
}

// Trigger implements parser.BlockParser.
func (p *BlockParser) Trigger() []byte {
	return []byte{fenceChar}
}

// Open implements parser.BlockParser.
func (p *BlockParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, segment := reader.PeekLine()
	pos := pc.BlockOffset()
	if pos < 0 {
		return nil, parser.NoChildren
	}

	fenceLength := countFence(line[pos:])
	if fenceLength < minFenceLength {
		return nil, parser.NoChildren
	}

	rest := line[pos+fenceLength:]

	nameLength := scanName(rest)
	if nameLength == 0 {
		return nil, parser.NoChildren
	}

	name := string(rest[:nameLength])
	rest = rest[nameLength:]

	var attributes parser.Attributes

	if !util.IsBlank(rest) {
		attrReader := text.NewReader(rest)

		attrs, ok := parser.ParseAttributes(attrReader)
		if !ok {
			return nil, parser.NoChildren
		}

		remaining, _ := attrReader.PeekLine()
		if !util.IsBlank(remaining) {
			return nil, parser.NoChildren
		}

		attributes = attrs
	}

	node := NewNode(name, fenceLength, reader.LineOffset()+pc.BlockIndent())

	for _, attr := range attributes {
		node.SetAttribute(attr.Name, directive.NormalizeAttributeValue(attr.Value))
	}

	advanceLine(reader, line, segment)

	return node, parser.HasChildren
}

// Continue implements parser.BlockParser.
func (p *BlockParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	container, ok := node.(*Node)
	if !ok {
		return parser.Close
	}

	line, segment := reader.PeekLine()

	column, ok := closingFenceColumn(line, reader.LineOffset(), container.fenceLength)
	if !ok {
		return parser.Continue | parser.HasChildren
	}

	// Let the innermost opened container consume the closing fence
	if inner := innermostOpenedContainer(container); inner != nil && column >= inner.column {
		return parser.Continue | parser.HasChildren
	}

	advanceLine(reader, line, segment)

	return parser.Close
}

// Close implements parser.BlockParser.
func (p *BlockParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	if container, ok := node.(*Node); ok {
		container.closed = true
	}
}

// CanInterruptParagraph implements parser.BlockParser.
func (p *BlockParser) CanInterruptParagraph() bool {
	return true
}

// CanAcceptIndentedLine implements parser.BlockParser.
func (p *BlockParser) CanAcceptIndentedLine() bool {
	return false
}

// closingFenceColumn returns the column of the closing fence found on the
// given line, if any
func closingFenceColumn(line []byte, lineOffset int, fenceLength int) (int, bool) {
	w, pos := util.IndentWidth(line, lineOffset)
	if w > 3 || pos >= len(line) {
		return 0, false
	}

	length := countFence(line[pos:])
	if length < fenceLength || !util.IsBlank(line[pos+length:]) {
		return 0, false
	}

	return lineOffset + w, true
}

func innermostOpenedContainer(node ast.Node) *Node {
	var innermost *Node

	for child := node.LastChild(); child != nil; child = child.LastChild() {
		if container, ok := child.(*Node); ok && !container.closed {
			innermost = container
		}
	}

	return innermost
}

func advanceLine(reader text.Reader, line []byte, segment text.Segment) {
	newline := 1
	if len(line) == 0 || line[len(line)-1] != '\n' {
		newline = 0
	}

	reader.Advance(segment.Stop - segment.Start - newline + segment.Padding)
}

func countFence(b []byte) int {
	i := 0
	for ; i < len(b) && b[i] == fenceChar; i++ {
	}
	return i
}

func scanName(b []byte) int {
	if len(b) == 0 || !isLetter(b[0]) {
		return 0
	}

	i := 1
	for ; i < len(b) && (isLetter(b[i]) || isDigit(b[i]) || b[i] == '-' || b[i] == '_'); i++ {
	}

	return i
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

var _ parser.BlockParser = &BlockParser{}
