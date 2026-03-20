package mcp

import (
	"context"
	"os"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/pkg/errors"
)

func serve(workspace string) error {
	root, err := os.OpenRoot(workspace)
	if err != nil {
		return errors.Wrapf(err, "could not open workspace %q", workspace)
	}
	defer root.Close()

	s := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    "amatl",
		Version: "1.0.0",
	}, nil)

	srv := &mcpServer{root: root}
	registerTools(s, srv)

	return s.Run(context.Background(), &sdkmcp.StdioTransport{})
}

func registerTools(s *sdkmcp.Server, srv *mcpServer) {
	sdkmcp.AddTool(s, toolDefListFiles(), srv.handleListFiles)
	sdkmcp.AddTool(s, toolDefTableOfContents(), srv.handleTableOfContents)
	sdkmcp.AddTool(s, toolDefReadSection(), srv.handleReadSection)
	sdkmcp.AddTool(s, toolDefFindSections(), srv.handleFindSections)
	sdkmcp.AddTool(s, toolDefListSections(), srv.handleListSections)
	sdkmcp.AddTool(s, toolDefUpdateSection(), srv.handleUpdateSection)
}

func toolDefListFiles() *sdkmcp.Tool {
	return &sdkmcp.Tool{
		Name: "list_files",
		Description: `List files available in the workspace.

Returns file paths relative to the workspace root. Use these paths directly with other tools (table_of_contents, read_section, find_sections, update_section).

The optional "pattern" argument filters files by name using shell glob syntax (* matches any sequence of non-separator characters). Defaults to "*.md" to list all Markdown files. Use "*" to list all files regardless of extension.`,
	}
}

func toolDefTableOfContents() *sdkmcp.Tool {
	return &sdkmcp.Tool{
		Name: "table_of_contents",
		Description: `Return the hierarchical table of contents of a Markdown file.

The "file" argument is a path relative to the workspace root (use list_files to discover available files).

The "entries" field is a JSON-encoded array. Each entry contains:
- "selector": a CSS-like selector that uniquely identifies the section (e.g. "h2#introduction"). Use it directly with read_section, find_sections, or update_section.
- "title": the heading text.
- "level": the heading level (1–6).
- "children": nested sub-sections (same structure, JSON-encoded).

RECOMMENDED FIRST STEP: always call table_of_contents before read_section or update_section so you have the correct selectors for the file.`,
	}
}

func toolDefReadSection() *sdkmcp.Tool {
	return &sdkmcp.Tool{
		Name: "read_section",
		Description: `Read the full content of a section in a Markdown file.

The "file" argument is a path relative to the workspace root (use list_files to discover available files).

A heading selector (e.g. "h2#introduction") returns the heading AND all its content (paragraphs, code blocks, sub-headings, etc.) until the next heading of equal or higher level. The returned "markdown" field contains the complete Markdown source of that section.

WORKFLOW for editing a section:
  1. Call table_of_contents to get the available selectors.
  2. Call read_section with the chosen selector to retrieve the current content.
  3. Modify the content in memory (keeping the heading line).
  4. Call update_section with the same selector and the full modified content.

SELECTOR SYNTAX:
  h2#my-section          — heading by ID (preferred, stable)
  h2:contains("Title")   — heading by text (fallback if no ID)
  code[lang="go"]        — fenced code block by language
  h2#api ~ table         — table following a specific heading`,
	}
}

func toolDefFindSections() *sdkmcp.Tool {
	return &sdkmcp.Tool{
		Name: "find_sections",
		Description: `Find all nodes matching a CSS-like selector anywhere in a Markdown file (full recursive search).

The "file" argument is a path relative to the workspace root (use list_files to discover available files).

Unlike read_section which searches only at the top level, find_sections walks the entire document tree. Useful for finding all code blocks of a given language, all blockquotes, tables, etc.

For heading selectors, each match is expanded to its full section (heading + content until the next same-or-higher-level heading).

SELECTOR SYNTAX:
  h2              — all level-2 headings (with their sections)
  code[lang="go"] — all Go code blocks
  blockquote p    — paragraphs inside blockquotes
  h2 ~ table      — tables that follow an h2 heading`,
	}
}

func toolDefListSections() *sdkmcp.Tool {
	return &sdkmcp.Tool{
		Name: "list_sections",
		Description: `List all headings in a Markdown file as a flat list with their CSS-like selectors.

The "file" argument is a path relative to the workspace root (use list_files to discover available files).

Each item contains "selector", "title", and "level". For a hierarchical view showing parent/child relationships, prefer table_of_contents instead.`,
	}
}

func toolDefUpdateSection() *sdkmcp.Tool {
	return &sdkmcp.Tool{
		Name: "update_section",
		Description: `Replace the entire content of a section in a Markdown file.

The "file" argument is a path relative to the workspace root (use list_files to discover available files).

IMPORTANT — a heading selector targets the FULL SECTION, not just the heading line:
  - "h2#introduction" replaces the h2 heading AND everything after it until the next h2 (or h1).
  - The "content" argument must include the heading line itself plus all the body you want to keep.
  - If "content" contains only the heading line, all body paragraphs will be erased.

REQUIRED WORKFLOW to avoid data loss:
  1. Call table_of_contents to identify the correct selector.
  2. Call read_section with that selector to retrieve the current full content.
  3. Modify the retrieved content (preserving unchanged parts).
  4. Call update_section with the same selector and the complete modified content.

SELECTOR SYNTAX: same as read_section (prefer ID-based selectors from table_of_contents).`,
	}
}
