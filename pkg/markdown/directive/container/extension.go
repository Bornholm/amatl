package container

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// BlockParserPriority is the priority of the container block parser.
//
// Goldmark built-in block parsers range from 100 (setext headings) to
// 1000 (paragraphs). The ':' trigger is not shared with any of them,
// so the value only has to stay below the paragraph parser (1000).
const BlockParserPriority = 750

// Extension registers the container directive parser and HTML renderer
type Extension struct {
}

// Extend implements goldmark.Extender.
func (e *Extension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(&BlockParser{}, BlockParserPriority),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&HTMLRenderer{}, 500),
		),
	)
}

var _ goldmark.Extender = &Extension{}
