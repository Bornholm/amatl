# Directives

> **Available for:** `Markdown`, `HTML`, `PDF`

`amatl` extends the CommonMark syntax with what we call "directives".

A directive takes this form:

```
:directive{attr="value1", attr="value2", ...}
```

Each directive triggers a specific behavior based on its type.

## Available directives

### `:include{url="<url>", fromHeadings="<headingLevel>", shiftHeadings="<levelShift>"}`

Include another Markdown file in your document.

#### Parameters

##### `url="<url>"`

**Required**

The URL of the Markdown document to include. This can be a local file or a remote document's URL (see ["URL resolving"](../url-resolving/README.md)).

##### `fromHeadings="<headingLevel>"`

**Optional**
**Type: `int`**

Only include sections that match the given heading level threshold.

##### `shiftHeadings="<levelShift>"`

**Optional**
**Type: `int`**

Shift the included headings by the given amount.

### `:toc{minLevel="<minLevel>", maxLevel="<maxLevel>"}`

Generate a table of contents for the whole document.

#### Parameters

##### `minLevel="<minLevel>"`

**Optional**
**Type: `int`**

Only include headings matching this minimum level.

##### `maxLevel="<maxLevel>"`

**Optional**
**Type: `int`**

Only include headings matching this maximum level.
