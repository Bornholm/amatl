# Directives

> **Available for:** `Markdown`, `HTML`, `PDF`

Amatl extends the CommonMark syntax with what we call "directives".

A directive takes this form:

```

:directive{attr="value1", attr="value2", ...}

```

Please take note of the linefeed before and after the directive. **They are required**.

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
