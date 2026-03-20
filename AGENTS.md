# AGENTS.md

This file provides guidance to LLM agents when working with code in this repository.

## Project

**amatl** is a CLI tool that transforms CommonMark (Markdown) into HTML, PDF, and processed Markdown. It uses a pipeline-based transformation architecture with middleware-composable transformers.

## Commands

```bash
make build              # Build binary to bin/amatl (CGO_ENABLED=0)
go test ./...           # Run all tests
go test ./pkg/resolver/... # Run tests for a specific package
make examples           # Generate example outputs (requires Chromium for PDF)
make website            # Build website from markdown sources
```

## Architecture

The core pipeline (`pkg/command/cli/render/pipeline.go`) chains these middlewares in order:

1. **Markdown** (`pkg/markdown/`) — Parses with Goldmark + extensions (GFM, syntax highlighting, Mermaid, frontmatter). Processes custom directives:
   - `include{}` — embeds other files, handles heading shifts and section filtering
   - `toc{}` — generates table of contents
   - `attrs{}` — applies HTML attributes to elements
2. **Template** — Applies Go templating (Sprig functions) with user-supplied vars
3. **HTML** — Renders markdown AST to HTML, applies a layout template
4. **PDF** (optional) — Drives headless Chromium via `chromedp` to produce PDF

Each middleware receives a `*pipeline.Payload` that carries the document content and contextual attributes.

### Resource Resolution (`pkg/resolver/`)

All file/URL lookups go through a scheme-based resolver (`file://`, `http://`, `https://`, `stdin://`, `amatl://`). The working directory is tracked in context and used to resolve relative paths. `amatl://` resolves against embedded built-in assets (layouts, etc.).

### HTML Layout System (`pkg/html/layout/`)

Built-in layouts (`amatl://document.html`, `amatl://presentation.html`, `amatl://website.html`) are Go templates. Custom layouts can be provided via URL. Layout variables are injected separately from document template vars.

### Transformer vs. Middleware

- `Transformer` is the core interface: `Transform(ctx context.Context, payload *Payload) error`
- `Middleware` wraps a transformer to compose chains
- Markdown directives are implemented as Goldmark AST node transformers that run during parsing

## Key Implementation Notes

- Path handling must work on both Unix and Windows. There is dedicated logic in `pkg/resolver/` to normalize paths between platforms — be careful when manipulating file paths.
- The `amatl://` scheme resolves against embedded FS assets (see `pkg/html/layout/` and related `embed` directives).
- Config files support relative paths that are rewritten to absolute based on the config file's location before processing.
- PDF generation requires a Chromium/Chrome binary. The `--pdf-no-sandbox` flag exists for environments where sandboxing is unavailable (e.g., CI).
