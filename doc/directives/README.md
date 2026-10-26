# Directives

> **Available for:** `Markdown`, `HTML`, `PDF`

Amatl extends the CommonMark syntax with what we call "directives".

A directive takes this form:

```

:directive{attr="value1", attr="value2", ...}

```

Please take note of the linefeed before and after the directive. **They are required**.

See also ["Container directives"](#container-directives) to wrap several blocks of content.

Each directive triggers a specific behavior based on its type.

## `:include{url="<url>", select="<selector>", fromHeadings="<headingLevel>", shiftHeadings="<levelShift>"}`

Include another Markdown file in your document.

### Parameters

#### `url="<url>"`

- **Required**

The URL of the Markdown document to include. This can be a local file or a remote document's URL (see ["URL resolving"](../url-resolving/README.md)).

#### `select="<selector>"`

- **Optional**
- **Type: `string`**

Only include the top-level nodes of the included document that match the given selector. Uses the same selector syntax as the MCP server (see ["MCP Server"](../mcp/README.md#-selector-syntax)).

For example, to include only the `## API` section of a file:

```
:include{url="./reference.md", select="h2#api"}
```

#### `fromHeadings="<headingLevel>"`

- **Optional**
- **Type: `int`**

Only include sections that match the given heading level threshold.

#### `shiftHeadings="<levelShift>"`

- **Optional**
- **Type: `int`**

Shift the included headings by the given amount.

## `:toc{minLevel="<minLevel>", maxLevel="<maxLevel>"}`

Generate a table of contents for the whole document.

### Parameters

#### `minLevel="<minLevel>"`

- **Optional**
- **Type: `int`**

Only include headings matching this minimum level.

#### `maxLevel="<maxLevel>"`

- **Optional**
- **Type: `int`**

Only include headings matching this maximum level.

## `:attrs{attributes...}`

Assign the given attributes to the following element (for example `class`, `id`, etc).

## Container directives

Container directives wrap several blocks of Markdown in an HTML element, which makes it possible to build layouts (columns, cards, callouts, sections...) without writing raw HTML.

A container is opened by a line of at least three colons followed by its name and optional attributes, and closed by a line of at least the same number of colons:

```
:::card{#security .highlight}
### Security

Sovereign hosting, end-to-end encryption.
:::
```

Containers can be nested. A closing fence always closes the innermost opened container, so the same fence length can be used at every level:

```
:::columns{count=3}
:::card
### Security
:::
:::card
### Simplicity
:::
:::card
### Sobriety
:::
:::
```

Other directives, such as `:include{}`, can be used inside containers.

### HTML rendering

- The element receives the `amatl-<name>` class, merged with the classes given in the attributes.
- The `article`, `aside`, `figcaption`, `figure`, `footer`, `header`, `main`, `nav` and `section` containers are rendered with the matching HTML element, other names with a `<div>`.
- Global HTML attributes (`id`, `class`, `style`, `title`, `role`...) as well as `data-*` and `aria-*` attributes are rendered as is. Other attributes are rendered as `data-*` attributes: `:::columns{count=3}` gives `<div class="amatl-columns" data-count="3">`.

An `:attrs{}` directive placed just before a container applies its attributes to the container.

### Built-in styles

The built-in layouts include the `amatl://components.css` stylesheet, which provides default styles for the following containers:

| Container  | Description                                                                         |
| ---------- | ----------------------------------------------------------------------------------- |
| `columns`  | Lays out its children in columns. The number of columns is set with `count` (1 to 6, defaults to 2). |
| `grid`     | Lays out its children in a responsive grid.                                         |
| `card`     | A bordered box.                                                                     |
| `callout`  | A highlighted box with an accent border.                                            |

These styles can be tuned with CSS custom properties (`--amatl-gap`, `--amatl-grid-min`, `--amatl-radius`, `--amatl-border-color`, `--amatl-card-bg`, `--amatl-accent`, `--amatl-callout-bg`). Containers with any other name get no default style and can be styled with your own stylesheet through their `amatl-<name>` class.
