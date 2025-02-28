package node

import (
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/yuin/goldmark/ast"
	extAST "github.com/yuin/goldmark/extension/ast"
)

func Renderers() map[ast.NodeKind]markdown.NodeRenderer {
	return map[ast.NodeKind]markdown.NodeRenderer{
		ast.KindDocument:         &DocumentRenderer{},
		ast.KindText:             &TextRenderer{},
		ast.KindString:           &StringRenderer{},
		ast.KindHeading:          WithLineSpacingBefore(&HeadingRenderer{}, 2),
		ast.KindAutoLink:         &AutoLinkRenderer{},
		ast.KindCodeBlock:        &CodeBlockRenderer{},
		ast.KindFencedCodeBlock:  WithLineSpacingBefore(&CodeBlockRenderer{}, 2),
		ast.KindHTMLBlock:        WithListOrHTMLBlockSpacing(&HTMLBlockRenderer{}),
		ast.KindRawHTML:          WithListOrHTMLBlockSpacing(&RawHTMLRenderer{}),
		ast.KindThematicBreak:    WithLineSpacingBefore(&ThematicBreakRenderer{}, 2),
		ast.KindListItem:         &ListItemRenderer{},
		ast.KindLink:             &LinkRenderer{},
		ast.KindImage:            &ImageRenderer{},
		ast.KindEmphasis:         &EmphasisRenderer{},
		ast.KindCodeSpan:         &CodeSpanRenderer{},
		ast.KindBlockquote:       WithLineSpacingBefore(&BlockquoteRenderer{}, 2),
		ast.KindParagraph:        WithLineSpacingBefore(&DummyRenderer{}, 2),
		ast.KindTextBlock:        WithLineSpacingBefore(&DummyRenderer{}, 2),
		ast.KindList:             WithListOrHTMLBlockSpacing(&DummyRenderer{}),
		extAST.KindTableCell:     &DummyRenderer{},
		extAST.KindStrikethrough: &StrikethroughRenderer{},
		extAST.KindTaskCheckBox:  &TaskCheckboxRenderer{},
		extAST.KindTable:         WithLineSpacingBefore(&TableRenderer{}, 2),
	}
}

func WithLineSpacingBefore(renderer markdown.NodeRenderer, count int) markdown.NodeRenderer {
	return markdown.WithBeforeRender(
		renderer,
		func(r *markdown.Render, node ast.Node, entering bool) {
			if !entering || node.PreviousSibling() == nil {
				return
			}

			for i := 0; i < count; i++ {
				_, _ = r.Writer().Write(markdown.NewLineChar)
			}
		},
	)
}

func WithListOrHTMLBlockSpacing(renderer markdown.NodeRenderer) markdown.NodeRenderer {
	return markdown.WithBeforeRender(
		renderer,
		func(r *markdown.Render, node ast.Node, entering bool) {
			if !entering || node.PreviousSibling() == nil {
				return
			}

			_, _ = r.Writer().Write(markdown.NewLineChar)

			if node.Type() != ast.TypeInline && node.HasBlankPreviousLines() {
				_, _ = r.Writer().Write(markdown.NewLineChar)
			}
		},
	)
}
