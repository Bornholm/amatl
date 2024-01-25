package markdown

import (
	"bytes"
	"go/format"
	"io"
	"strconv"
	"unicode/utf8"
	"unsafe"

	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
	extAST "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
)

var (
	NewLineChar             = []byte{'\n'}
	SpaceChar               = []byte{' '}
	StrikeThroughChars      = []byte("~~")
	ThematicBreakChars      = []byte("---")
	BlockquoteChars         = []byte{'>', ' '}
	CodeBlockChars          = []byte("```")
	TableHeaderColChar      = []byte{'-'}
	TableHeaderAlignColChar = []byte{':'}
	Heading1UnderlineChar   = []byte{'='}
	Heading2UnderlineChar   = []byte{'-'}
	FourSpacesChars         = bytes.Repeat([]byte{' '}, 4)
)

// Ensure compatibility with Goldmark parser.
var _ renderer.Renderer = &Renderer{}

// Renderer allows to render markdown AST into markdown bytes in consistent format.
// Render is reusable across Renders, it holds configuration only.
type Renderer struct {
	underlineHeadings bool
	softWraps         bool
	emphToken         []byte
	strongToken       []byte // if nil, use emphToken*2
	listIndentStyle   ListIndentStyle

	// language name => format function
	formatters map[string]func([]byte) []byte
	renderers  map[ast.NodeKind]NodeRenderer
}

type NodeRenderer interface {
	Render(r *Render, node ast.Node, entering bool) (ast.WalkStatus, error)
}

// AddOptions pulls Markdown renderer specific options from the given list,
// and applies them to the renderer.
func (mr *Renderer) AddOptions(opts ...renderer.Option) {
	mdopts := make([]Option, 0, len(opts))
	for _, o := range opts {
		if mo, ok := o.(Option); ok {
			mdopts = append(mdopts, mo)
		}
	}
	mr.AddMarkdownOptions(mdopts...)
}

// AddMarkdownOptions modifies the Renderer with the given options.
func (mr *Renderer) AddMarkdownOptions(opts ...Option) {
	for _, o := range opts {
		o.apply(mr)
	}
}

// Option customizes the behavior of the markdown renderer.
type Option interface {
	renderer.Option

	apply(r *Renderer)
}

type optionFunc func(*Renderer)

func (f optionFunc) SetConfig(*renderer.Config) {}

func (f optionFunc) apply(r *Renderer) {
	f(r)
}

// WithNodeRenderer configures the renderer to use a custom renderer for the nodes of the given kind.
func WithNodeRenderer(kind ast.NodeKind, renderer NodeRenderer) Option {
	return optionFunc(func(r *Renderer) {
		if r.renderers == nil {
			r.renderers = make(map[ast.NodeKind]NodeRenderer)
		}

		r.renderers[kind] = renderer
	})
}

// WithUnderlineHeadings configures the renderer to use
// Setext-style headers (=== and ---).
func WithUnderlineHeadings() Option {
	return optionFunc(func(r *Renderer) {
		r.underlineHeadings = true
	})
}

// WithSoftWraps allows you to wrap lines even on soft line breaks.
func WithSoftWraps() Option {
	return optionFunc(func(r *Renderer) {
		r.softWraps = true
	})
}

// WithEmphasisToken specifies the character used to wrap emphasised text.
// Per the CommonMark spec, valid values are '*' and '_'.
//
// Defaults to '*'.
func WithEmphasisToken(c rune) Option {
	return optionFunc(func(r *Renderer) {
		buf := make([]byte, 4) // enough to encode any utf8 rune
		n := utf8.EncodeRune(buf, c)
		r.emphToken = buf[:n]
	})
}

// WithStrongToken specifies the string used to wrap bold text.
// Per the CommonMark spec, valid values are '**' and '__'.
//
// Defaults to repeating the emphasis token twice.
// See [WithEmphasisToken] for how to change that.
func WithStrongToken(s string) Option {
	return optionFunc(func(r *Renderer) {
		r.strongToken = []byte(s)
	})
}

// ListIndentStyle specifies how items nested inside lists
// should be indented.
type ListIndentStyle int

const (
	// ListIndentAligned specifies that items inside a list item
	// should be aligned to the content in the first item.
	//
	//	- First paragraph.
	//
	//	  Second paragraph aligned with the first.
	//
	// This applies to ordered lists too.
	//
	//	1. First paragraph.
	//
	//	   Second paragraph aligned with the first.
	//
	//	...
	//
	//	10. Contents.
	//
	//	    Long lists indent content further.
	//
	// This is the default.
	ListIndentAligned ListIndentStyle = iota

	// ListIndentUniform specifies that items inside a list item
	// should be aligned uniformly with 4 spaces.
	//
	// For example:
	//
	//	- First paragraph.
	//
	//	    Second paragraph indented 4 spaces.
	//
	// For ordered lists:
	//
	//	1. First paragraph.
	//
	//	    Second paragraph indented 4 spaces.
	//
	//	...
	//
	//	10. Contents.
	//
	//	    Always indented 4 spaces.
	ListIndentUniform
)

// WithListIndentStyle specifies how contents nested under a list item
// should be indented.
//
// Defaults to [ListIndentAligned].
func WithListIndentStyle(style ListIndentStyle) Option {
	return optionFunc(func(r *Renderer) {
		r.listIndentStyle = style
	})
}

// CodeFormatter reformats code samples found in the document,
// matching them by name.
type CodeFormatter struct {
	// Name of the language.
	Name string

	// Aliases for the language, if any.
	Aliases []string

	// Function to format the code snippet.
	// In case of errors, format functions should typically return
	// the original string unchanged.
	Format func([]byte) []byte
}

// GoCodeFormatter is a [CodeFormatter] that reformats Go source code inside
// fenced code blocks tagged with 'go' or 'Go'.
//
//	```go
//	func main() {
//	}
//	```
//
// Supply it to the renderer with [WithCodeFormatters].
var GoCodeFormatter = CodeFormatter{
	Name:    "go",
	Aliases: []string{"Go"},
	Format:  formatGo,
}

func formatGo(src []byte) []byte {
	gofmt, err := format.Source(src)
	if err != nil {
		// We don't handle gofmt errors.
		// If code is not compilable we just
		// don't format it without any warning.
		return src
	}
	return gofmt
}

// WithCodeFormatters changes the functions used to reformat code blocks found
// in the original file.
//
//	formatters := []markdown.CodeFormatter{
//		markdown.GoCodeFormatter,
//		// ...
//	}
//	r := NewRenderer()
//	r.AddMarkdownOptions(WithCodeFormatters(formatters...))
//
// Defaults to empty.
func WithCodeFormatters(fs ...CodeFormatter) Option {
	return optionFunc(func(r *Renderer) {
		formatters := make(map[string]func([]byte) []byte, len(fs))
		for _, f := range fs {
			formatters[f.Name] = f.Format
			for _, alias := range f.Aliases {
				formatters[alias] = f.Format
			}
		}
		r.formatters = formatters
	})
}

// NewRenderer builds a new Markdown renderer with default settings.
// To use this with goldmark.Markdown, use the goldmark.WithRenderer option.
//
//	r := markdown.NewRenderer()
//	md := goldmark.New(goldmark.WithRenderer(r))
//	md.Convert(src, w)
//
// Alternatively, you can call [Renderer.Render] directly.
//
//	r := markdown.NewRenderer()
//	r.Render(w, src, node)
//
// Use [Renderer.AddMarkdownOptions] to customize the output of the renderer.
func NewRenderer() *Renderer {
	return &Renderer{
		emphToken: []byte{'*'},
		// Leave strongToken as nil by default.
		// At render time, we'll use what was specified,
		// or repeat emphToken twice to get the strong token.
		renderers: DefaultRenderers(),
	}
}

func (mr *Renderer) newRender(w io.Writer, source []byte) *Render {
	strongToken := mr.strongToken
	if len(strongToken) == 0 {
		strongToken = bytes.Repeat(mr.emphToken, 2)
	}

	return &Render{
		mr:          mr,
		w:           wrapWithLineIndentWriter(w),
		source:      source,
		strongToken: strongToken,
		emphToken:   mr.emphToken,
	}
}

// Render renders the given AST node to the given writer,
// given the original source from which the node was parsed.
//
// NOTE: This is the entry point used by Goldmark.
func (mr *Renderer) Render(w io.Writer, source []byte, node ast.Node) error {
	// Perform DFS.
	return ast.Walk(node, mr.newRender(w, source).renderNode)
}

func (r *Render) renderNode(node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering && node.PreviousSibling() != nil {
		switch node.(type) {
		// All Block types (except few) usually have 2x new lines before itself when they are non-first siblings.
		case *ast.Paragraph, *ast.Heading, *ast.FencedCodeBlock,
			*ast.CodeBlock, *ast.ThematicBreak, *extAST.Table,
			*ast.Blockquote:
			_, _ = r.w.Write(NewLineChar)
			_, _ = r.w.Write(NewLineChar)
		case *ast.List, *ast.HTMLBlock:
			_, _ = r.w.Write(NewLineChar)
			if node.HasBlankPreviousLines() {
				_, _ = r.w.Write(NewLineChar)
			}
		case *ast.ListItem:
			// TODO(bwplotka): Handle tight/loose rule explicitly.
			// See: https://github.github.com/gfm/#loose
			if node.HasBlankPreviousLines() {
				_, _ = r.w.Write(NewLineChar)
			}
		}
	}

	kind := node.Kind()
	renderer, exists := r.mr.renderers[kind]
	if !exists {
		return ast.WalkStop, errors.Errorf("no renderer registered for kind '%s'", kind)
	}

	return renderer.Render(r, node, entering)
}

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
	return []byte{parList.Marker, SpaceChar[0]}
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
				if _, err := w.Write(SpaceChar); err != nil {
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
