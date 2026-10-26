# Directives

> **Available for:** `Markdown`, `HTML`, `PDF`

Amatl extends the CommonMark syntax with what we call "directives".

A directive takes this form:

```

:directive{attr="value1", attr="value2", ...}

```

Please take note of the linefeed before and after the directive. **They are required**.

An attribute value is either a quoted string or a bare number or boolean: `minLevel="2"` and `minLevel=2` mean the same thing.

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

A container directive wraps several Markdown blocks in one HTML element. Use it to lay out columns, cards or sections without writing HTML.

A line of three or more colons, followed by a name and optional attributes, opens a container. A line of at least as many colons closes it:

```
:::card{#security .highlight}
### Security

Sovereign hosting, end-to-end encryption.
:::
```

Containers nest. A closing fence closes the innermost open container, so you can use three colons at every level:

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

Indentation matters when a container sits in a list item. A closing fence closes the deepest open container whose opening fence is not more indented than the closing fence. Everything opened inside that container closes with it. In practice, indent a closing fence like its opening fence.

Other directives work inside a container, `:include{}` included.

### HTML output

The container becomes an element with the `amatl-<name>` class, plus the classes given in its attributes.

The `article`, `aside`, `figcaption`, `figure`, `footer`, `header`, `main`, `nav` and `section` containers use the HTML element of the same name. Any other name gives a `<div>`.

Global HTML attributes such as `id`, `class`, `style`, `title` or `role` are copied as they are, and so are `data-*` and `aria-*` attributes. Any other attribute becomes a `data-*` attribute, except `bg` (see [Background image](#background-image)). For example, `:::columns{count=3}` gives `<div class="amatl-columns" data-count="3">`.

An `:attrs{}` directive placed right before a container adds its attributes to the container.

### Background image

The `bg` attribute sets a background image on a container:

```
:::section{.hero bg="./cover.jpg"}
# Our offer
:::
```

Amatl resolves the path like an image path, so a relative path in an included file points next to that file. HTML and PDF output embed the image as a `data:` URL and pass it to CSS through the `--amatl-bg` custom property:

```html
<section class="amatl-section hero" style="--amatl-bg: url(&quot;data:image/jpeg;base64,...&quot;)">
```

`amatl://components.css` turns it into a centered background that covers the element. A custom layout can use `var(--amatl-bg)` in its own rules. Markdown output keeps the path. In an included file, amatl makes it absolute, as it does for images.

### Built-in styles

The built-in layouts load `amatl://components.css`, which styles four containers:

| Container | Rendering                                                                 |
| --------- | ------------------------------------------------------------------------- |
| `columns` | one column per child, 2 by default, set `count` from 1 to 6 to change it  |
| `grid`    | as many columns as fit the page width                                     |
| `card`    | a box with a border and rounded corners                                   |
| `callout` | a tinted box with a colored left border                                   |

In `columns` and `grid`, a horizontal rule (`---`) spans the whole row instead of taking a cell.

Change the spacing and colors with these CSS custom properties: `--amatl-gap`, `--amatl-grid-min`, `--amatl-radius`, `--amatl-border-color`, `--amatl-card-bg`, `--amatl-accent` and `--amatl-callout-bg`.

A container with any other name has no default style. Style it in your own layout through its `amatl-<name>` class.
