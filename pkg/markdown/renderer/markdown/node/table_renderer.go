package node

import (
	"bytes"
	"fmt"

	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/mattn/go-runewidth"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
	extAST "github.com/yuin/goldmark/extension/ast"
)

type TableRenderer struct {
}

// Render implements NodeRenderer.
func (tr *TableRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	table, ok := node.(*extAST.Table)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *extAST.Table, got '%T'", node)
	}

	// Render it straight away. No nested tables are supported and we expect
	// tables to have limited content, so limit WALK.
	if err := tr.renderTable(r, table); err != nil {
		return ast.WalkStop, fmt.Errorf("rendering table: %w", err)
	}

	return ast.WalkSkipChildren, nil
}

func (tr *TableRenderer) renderTable(r *markdown.Render, node *extAST.Table) error {
	var (
		columnAligns []extAST.Alignment
		columnWidths []int
		colIndex     int
		cellBuf      bytes.Buffer
	)

	// Walk tree initially to count column widths and alignments.
	for n := node.FirstChild(); n != nil; n = n.NextSibling() {
		if err := ast.Walk(n, func(inner ast.Node, entering bool) (ast.WalkStatus, error) {
			switch tnode := inner.(type) {
			case *extAST.TableRow, *extAST.TableHeader:
				if entering {
					colIndex = 0
				}
			case *extAST.TableCell:
				if entering {
					if _, isHeader := tnode.Parent().(*extAST.TableHeader); isHeader {
						columnAligns = append(columnAligns, tnode.Alignment)
					}

					cellBuf.Reset()
					if err := ast.Walk(tnode, r.Renderer().NewRender(&cellBuf, r.Source()).RenderNode); err != nil {
						return ast.WalkStop, err
					}
					width := runewidth.StringWidth(cellBuf.String())
					if len(columnWidths) <= colIndex {
						columnWidths = append(columnWidths, width)
					} else if width > columnWidths[colIndex] {
						columnWidths[colIndex] = width
					}
					colIndex++
					return ast.WalkSkipChildren, nil
				}
			default:
				return ast.WalkStop, fmt.Errorf("detected unexpected tree type %v", tnode.Kind())
			}
			return ast.WalkContinue, nil
		}); err != nil {
			return err
		}
	}

	// Write all according to alignments and width.
	for n := node.FirstChild(); n != nil; n = n.NextSibling() {
		if err := ast.Walk(n, func(inner ast.Node, entering bool) (ast.WalkStatus, error) {
			switch tnode := inner.(type) {
			case *extAST.TableRow:
				if entering {
					colIndex = 0
					_, _ = r.Writer().Write(markdown.NewLineChar)
					break
				}

				_, _ = r.Writer().Write([]byte("|"))
			case *extAST.TableHeader:
				if entering {
					colIndex = 0
					break
				}

				_, _ = r.Writer().Write([]byte("|\n"))
				for i, align := range columnAligns {
					_, _ = r.Writer().Write([]byte{'|'})
					width := columnWidths[i]

					left, right := markdown.TableHeaderColChar, markdown.TableHeaderColChar
					switch align {
					case extAST.AlignLeft:
						left = markdown.TableHeaderAlignColChar
					case extAST.AlignRight:
						right = markdown.TableHeaderAlignColChar
					case extAST.AlignCenter:
						left, right = markdown.TableHeaderAlignColChar, markdown.TableHeaderAlignColChar
					}
					_, _ = r.Writer().Write(left)
					_, _ = r.Writer().Write(bytes.Repeat(markdown.TableHeaderColChar, width))
					_, _ = r.Writer().Write(right)
				}
				_, _ = r.Writer().Write([]byte("|"))
			case *extAST.TableCell:
				if !entering {
					break
				}

				width := columnWidths[colIndex]
				align := columnAligns[colIndex]

				if tnode.Parent().Kind() == extAST.KindTableHeader {
					align = extAST.AlignLeft
				}

				cellBuf.Reset()
				if err := ast.Walk(tnode, r.Renderer().NewRender(&cellBuf, r.Source()).RenderNode); err != nil {
					return ast.WalkStop, err
				}

				_, _ = r.Writer().Write([]byte("| "))
				whitespaceWidth := width - runewidth.StringWidth(cellBuf.String())
				switch align {
				default:
					fallthrough
				case extAST.AlignLeft:
					_, _ = r.Writer().Write(cellBuf.Bytes())
					_, _ = r.Writer().Write(bytes.Repeat([]byte{' '}, 1+whitespaceWidth))
				case extAST.AlignCenter:
					first := whitespaceWidth / 2
					_, _ = r.Writer().Write(bytes.Repeat([]byte{' '}, first))
					_, _ = r.Writer().Write(cellBuf.Bytes())
					_, _ = r.Writer().Write(bytes.Repeat([]byte{' '}, whitespaceWidth-first))
					_, _ = r.Writer().Write([]byte{' '})
				case extAST.AlignRight:
					_, _ = r.Writer().Write(bytes.Repeat([]byte{' '}, whitespaceWidth))
					_, _ = r.Writer().Write(cellBuf.Bytes())
					_, _ = r.Writer().Write([]byte{' '})
				}
				colIndex++
				return ast.WalkSkipChildren, nil
			default:
				return ast.WalkStop, fmt.Errorf("detected unexpected tree type %v", tnode.Kind())
			}
			return ast.WalkContinue, nil
		}); err != nil {
			return err
		}
	}
	return nil
}

var _ markdown.NodeRenderer = &TableRenderer{}
