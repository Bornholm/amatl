package markdown

import (
	"github.com/yuin/goldmark/ast"
	extAST "github.com/yuin/goldmark/extension/ast"
)

func DefaultRenderers() map[ast.NodeKind]NodeRenderer {
	return map[ast.NodeKind]NodeRenderer{
		ast.KindDocument:         &DocumentRenderer{},
		ast.KindText:             &TextRenderer{},
		ast.KindString:           &StringRenderer{},
		ast.KindHeading:          WithLineSpacingBefore(&HeadingRenderer{}),
		ast.KindAutoLink:         &AutoLinkRenderer{},
		ast.KindCodeBlock:        &CodeBlockRenderer{},
		ast.KindFencedCodeBlock:  WithLineSpacingBefore(&CodeBlockRenderer{}),
		ast.KindHTMLBlock:        WithLineSpacingBefore(&HTMLBlockRenderer{}),
		ast.KindThematicBreak:    WithLineSpacingBefore(&ThematicBreakRenderer{}),
		ast.KindListItem:         &ListItemRenderer{},
		ast.KindLink:             &LinkRenderer{},
		ast.KindImage:            &ImageRenderer{},
		ast.KindEmphasis:         &EmphasisRenderer{},
		ast.KindCodeSpan:         &CodeSpanRenderer{},
		ast.KindBlockquote:       WithLineSpacingBefore(&BlockquoteRenderer{}),
		ast.KindParagraph:        WithLineSpacingBefore(&DummyRenderer{}),
		ast.KindTextBlock:        WithLineSpacingBefore(&DummyRenderer{}),
		ast.KindList:             WithLineSpacingBefore(&DummyRenderer{}),
		extAST.KindTableCell:     &DummyRenderer{},
		extAST.KindStrikethrough: &StrikethroughRenderer{},
		extAST.KindTaskCheckBox:  &TaskCheckboxRenderer{},
		extAST.KindTable:         WithLineSpacingBefore(&TableRenderer{}),
	}
}

func WithLineSpacingBefore(renderer NodeRenderer) NodeRenderer {
	return WithBeforeRender(
		renderer,
		func(r *Render, node ast.Node, entering bool) {
			if !entering || node.PreviousSibling() == nil {
				return
			}

			if node.HasBlankPreviousLines() {
				return
			}

			_, _ = r.w.Write(NewLineChar)
		},
	)
}
