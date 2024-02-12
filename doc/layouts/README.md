# Layouts

> **Available for:** `HTML`, `PDF`

A layout is a HTML template that can be used to "decorates" your Markdown content.
It permits multiples things:

- Insert content in front or after your Markdown content;
- Style and transform your Markdown content to match your organization preferences;
- Use layout-specific variables to transform/enhance its content.

`amatl` provides 2 base layouts:

- `amatl://document.html` - A layout to generate generic document in A4 format (used by default);
- `amatl://presentation.html` - A presentation ("slides") layout;

You can also use a custom layout passing the flag `--html-layout` to your command.

## `amatl://document.html`

### Example

[See `../../examples/document.pdf`](../../examples/document/document.pdf)

### Variables

#### `Title`

The `Title` variable is used to defined the HTML document title.

## `amatl://presentation.html`

### Example

[See `../../examples/presentation.pdf`](../../examples/presentation/presentation.pdf)

### Variables

#### `Title`

The `Title` variable is used to defined the HTML document title.