package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path"
	"strings"

	"github.com/Bornholm/amatl/pkg/markdown/selector"
	mdrenderer "github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown/node"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	goldmarkparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// mcpServer holds the workspace root and provides the MCP tool handlers.
// File access is confined to the workspace by os.Root — the OS kernel
// enforces the boundary, preventing any path-traversal escape.
type mcpServer struct {
	root *os.Root
}

// ---------------------------------------------------------------------------
// Input / output types
// ---------------------------------------------------------------------------

type listFilesInput struct {
	Pattern string `json:"pattern,omitempty"`
}

type listFilesOutput struct {
	Workspace string   `json:"workspace"`
	Files     []string `json:"files"`
}

type fileInput struct {
	File string `json:"file"`
}

type selectorInput struct {
	File     string `json:"file"`
	Selector string `json:"selector"`
}

type updateInput struct {
	File     string `json:"file"`
	Selector string `json:"selector"`
	Content  string `json:"content"`
}

// SectionInfo describes a heading found in a Markdown file.
type SectionInfo struct {
	Selector string `json:"selector"`
	Title    string `json:"title"`
	Level    int    `json:"level,omitempty"`
}

// TOCEntry is a node in the hierarchical table of contents.
type TOCEntry struct {
	Selector string      `json:"selector"`
	Title    string      `json:"title"`
	Level    int         `json:"level"`
	Children []*TOCEntry `json:"children,omitempty"`
	parent   *TOCEntry
}

type tocOutput struct {
	File    string `json:"file"`
	Entries string `json:"entries"`
}

type listSectionsOutput struct {
	Sections []SectionInfo `json:"sections"`
}

type markdownOutput struct {
	Markdown string `json:"markdown"`
}

type updateOutput struct {
	OK bool `json:"ok"`
}

// ---------------------------------------------------------------------------
// Tool handlers
// ---------------------------------------------------------------------------

func (srv *mcpServer) handleListFiles(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	args listFilesInput,
) (*sdkmcp.CallToolResult, listFilesOutput, error) {
	pattern := args.Pattern
	if pattern == "" {
		pattern = "*.md"
	}

	var files []string
	err := fs.WalkDir(srv.root.FS(), ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") && p != "." {
				return fs.SkipDir
			}
			return nil
		}
		matched, matchErr := path.Match(pattern, d.Name())
		if matchErr != nil {
			return matchErr
		}
		if matched {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return nil, listFilesOutput{}, errors.Wrap(err, "could not list workspace files")
	}

	return nil, listFilesOutput{Workspace: srv.root.Name(), Files: files}, nil
}

func (srv *mcpServer) handleTableOfContents(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	args fileInput,
) (*sdkmcp.CallToolResult, tocOutput, error) {
	doc, source, err := parseMarkdownFile(srv.root, args.File)
	if err != nil {
		return nil, tocOutput{}, errors.Wrapf(err, "could not parse file %q", args.File)
	}

	entries, err := buildTOCTree(doc, source)
	if err != nil {
		return nil, tocOutput{}, errors.Wrap(err, "could not build table of contents")
	}

	entriesJSON, err := json.Marshal(entries)
	if err != nil {
		return nil, tocOutput{}, errors.Wrap(err, "could not serialize table of contents")
	}

	return nil, tocOutput{File: args.File, Entries: string(entriesJSON)}, nil
}

func (srv *mcpServer) handleListSections(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	args fileInput,
) (*sdkmcp.CallToolResult, listSectionsOutput, error) {
	doc, source, err := parseMarkdownFile(srv.root, args.File)
	if err != nil {
		return nil, listSectionsOutput{}, errors.Wrapf(err, "could not parse file %q", args.File)
	}

	var sections []SectionInfo

	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		heading, ok := n.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}
		sections = append(sections, SectionInfo{
			Selector: headingSelector(heading, source),
			Title:    string(heading.Text(source)),
			Level:    heading.Level,
		})
		return ast.WalkContinue, nil
	})

	return nil, listSectionsOutput{Sections: sections}, nil
}

func (srv *mcpServer) handleReadSection(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	args selectorInput,
) (*sdkmcp.CallToolResult, markdownOutput, error) {
	doc, source, err := parseMarkdownFile(srv.root, args.File)
	if err != nil {
		return nil, markdownOutput{}, errors.Wrapf(err, "could not parse file %q", args.File)
	}

	sel, err := selector.Parse(args.Selector)
	if err != nil {
		return nil, markdownOutput{}, fmt.Errorf(
			"invalid selector %q: %w — use selectors returned by table_of_contents, or the syntax: h2#id, h2:contains(\"text\"), code[lang=\"go\"]",
			args.Selector, err,
		)
	}

	nodes := sel.MatchTopLevel(doc, source)
	if len(nodes) == 0 {
		return nil, markdownOutput{}, fmt.Errorf(
			"selector %q matched no section in %q — call table_of_contents on this file to get the list of valid selectors",
			args.Selector, args.File,
		)
	}

	md, err := renderNodesToMarkdown(nodes, source)
	if err != nil {
		return nil, markdownOutput{}, errors.Wrap(err, "could not render section to markdown")
	}

	return nil, markdownOutput{Markdown: md}, nil
}

func (srv *mcpServer) handleFindSections(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	args selectorInput,
) (*sdkmcp.CallToolResult, markdownOutput, error) {
	doc, source, err := parseMarkdownFile(srv.root, args.File)
	if err != nil {
		return nil, markdownOutput{}, errors.Wrapf(err, "could not parse file %q", args.File)
	}

	sel, err := selector.Parse(args.Selector)
	if err != nil {
		return nil, markdownOutput{}, fmt.Errorf(
			"invalid selector %q: %w — use selectors returned by table_of_contents, or the syntax: h2#id, h2:contains(\"text\"), code[lang=\"go\"]",
			args.Selector, err,
		)
	}

	nodes := sel.FindAll(doc, source)
	if len(nodes) == 0 {
		return nil, markdownOutput{}, fmt.Errorf(
			"selector %q matched nothing in %q — call table_of_contents to browse available sections, or find_sections with a broader selector (e.g. \"h2\", \"code\")",
			args.Selector, args.File,
		)
	}

	md, err := renderNodesToMarkdown(nodes, source)
	if err != nil {
		return nil, markdownOutput{}, errors.Wrap(err, "could not render found sections to markdown")
	}

	return nil, markdownOutput{Markdown: md}, nil
}

func (srv *mcpServer) handleUpdateSection(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	args updateInput,
) (*sdkmcp.CallToolResult, updateOutput, error) {
	doc, source, err := parseMarkdownFile(srv.root, args.File)
	if err != nil {
		return nil, updateOutput{}, errors.Wrapf(err, "could not parse file %q", args.File)
	}

	sel, err := selector.Parse(args.Selector)
	if err != nil {
		return nil, updateOutput{}, fmt.Errorf(
			"invalid selector %q: %w — use selectors returned by table_of_contents, or the syntax: h2#id, h2:contains(\"text\"), code[lang=\"go\"]",
			args.Selector, err,
		)
	}

	nodes := sel.MatchTopLevel(doc, source)
	if len(nodes) == 0 {
		return nil, updateOutput{}, fmt.Errorf(
			"selector %q matched no section in %q — call table_of_contents on this file to get the list of valid selectors, then call read_section to retrieve the current content before calling update_section",
			args.Selector, args.File,
		)
	}

	start, end := nodeSourceRange(nodes, source)
	if start < 0 || end <= start {
		return nil, updateOutput{}, fmt.Errorf(
			"could not determine the byte range for selector %q in %q — the matched nodes may not be block-level elements; use a heading selector (e.g. h2#my-section) to target a full section",
			args.Selector, args.File,
		)
	}

	content := args.Content
	if expectedLevel := headingLevelFromSelector(args.Selector); expectedLevel > 0 {
		content = normalizeFirstHeadingLevel(content, expectedLevel)
	}

	newContent := []byte(content)
	if len(newContent) > 0 && newContent[len(newContent)-1] != '\n' {
		newContent = append(newContent, '\n')
	}

	newSource := make([]byte, 0, len(source)-end+start+len(newContent))
	newSource = append(newSource, source[:start]...)
	newSource = append(newSource, newContent...)
	newSource = append(newSource, source[end:]...)

	if err := srv.root.WriteFile(args.File, newSource, 0o644); err != nil {
		return nil, updateOutput{}, errors.Wrapf(err, "could not write file %q", args.File)
	}

	return nil, updateOutput{OK: true}, nil
}

// ---------------------------------------------------------------------------
// TOC helpers
// ---------------------------------------------------------------------------

func buildTOCTree(doc ast.Node, source []byte) ([]*TOCEntry, error) {
	var roots []*TOCEntry
	var last *TOCEntry

	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		heading, ok := n.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}

		entry := &TOCEntry{
			Level:    heading.Level,
			Title:    string(heading.Text(source)),
			Selector: headingSelector(heading, source),
		}

		if last == nil {
			roots = append(roots, entry)
		} else if entry.Level > last.Level {
			entry.parent = last
			last.Children = append(last.Children, entry)
		} else {
			ancestor := last.parent
			for ancestor != nil && ancestor.Level >= entry.Level {
				ancestor = ancestor.parent
			}
			if ancestor == nil {
				roots = append(roots, entry)
			} else {
				entry.parent = ancestor
				ancestor.Children = append(ancestor.Children, entry)
			}
		}

		last = entry
		return ast.WalkContinue, nil
	})

	return roots, nil
}

// headingSelector builds the best CSS-like selector for a heading.
// Prefers ID-based selectors; falls back to :contains() if no ID is set.
func headingSelector(heading *ast.Heading, source []byte) string {
	tag := fmt.Sprintf("h%d", heading.Level)

	if id, exists := heading.AttributeString("id"); exists {
		var idStr string
		switch v := id.(type) {
		case []byte:
			idStr = string(v)
		case string:
			idStr = v
		}
		if idStr != "" {
			return tag + "#" + idStr
		}
	}

	title := string(heading.Text(source))
	return tag + `:contains("` + escapeContains(title) + `")`
}

// escapeContains escapes double quotes and backslashes inside a :contains() pattern.
func escapeContains(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}

// ---------------------------------------------------------------------------
// Markdown parsing and rendering helpers
// ---------------------------------------------------------------------------

func parseMarkdownFile(root *os.Root, name string) (ast.Node, []byte, error) {
	data, err := root.ReadFile(name)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "could not read file %q", name)
	}

	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	p := md.Parser()
	p.AddOptions(goldmarkparser.WithAutoHeadingID())

	doc := p.Parse(text.NewReader(data))
	return doc, data, nil
}

func renderNodesToMarkdown(nodes []ast.Node, source []byte) (string, error) {
	mr := mdrenderer.NewRenderer()
	mr.AddMarkdownOptions(mdrenderer.WithNodeRenderers(node.Renderers()))

	var buf bytes.Buffer
	for i, n := range nodes {
		if i > 0 {
			buf.WriteString("\n")
		}
		if err := mr.Render(&buf, source, n); err != nil {
			return "", errors.Wrapf(err, "could not render node of kind %q", n.Kind().String())
		}
	}
	return buf.String(), nil
}

// nodeSourceRange returns the [start, end) byte range of the given nodes.
// ATX headings store only the content (after '## ') in their Lines(), so
// we walk back to the true start of the line.
func nodeSourceRange(nodes []ast.Node, source []byte) (start, end int) {
	start = math.MaxInt
	end = 0
	for _, n := range nodes {
		lines := n.Lines()
		if lines == nil {
			continue
		}
		for i := 0; i < lines.Len(); i++ {
			seg := lines.At(i)
			if seg.Start < start {
				start = seg.Start
			}
			if seg.Stop > end {
				end = seg.Stop
			}
		}
	}
	if start == math.MaxInt {
		return -1, -1
	}
	// Walk back to the beginning of the line (ATX headings store content after '## ')
	for start > 0 && source[start-1] != '\n' {
		start--
	}
	// Advance end to include the trailing newline if not already included
	if end > 0 && end <= len(source) && source[end-1] != '\n' {
		for end < len(source) && source[end] != '\n' {
			end++
		}
		if end < len(source) {
			end++ // include the '\n'
		}
	}
	return
}

// headingLevelFromSelector extracts the heading level from a selector like
// "h2#id", "h3:contains(...)", etc. Returns 0 if the selector does not target
// a heading element.
func headingLevelFromSelector(sel string) int {
	if len(sel) < 2 || sel[0] != 'h' {
		return 0
	}
	c := sel[1]
	if c < '1' || c > '6' {
		return 0
	}
	return int(c - '0')
}

// normalizeFirstHeadingLevel finds the first ATX heading line in content and
// rewrites its level to match wantedLevel, leaving the title text unchanged.
func normalizeFirstHeadingLevel(content string, wantedLevel int) string {
	lines := strings.SplitAfter(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimRight(line, "\n")
		if !strings.HasPrefix(trimmed, "#") {
			continue
		}
		rest := strings.TrimLeft(trimmed, "#")
		currentLevel := len(trimmed) - len(rest)
		if currentLevel == wantedLevel {
			break
		}
		lines[i] = strings.Repeat("#", wantedLevel) + rest + "\n"
		break
	}
	return strings.Join(lines, "")
}
