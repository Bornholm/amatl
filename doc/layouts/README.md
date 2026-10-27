# Layouts

> **Available for:** `HTML`, `PDF`

A layout is the HTML template that wraps your rendered Markdown. Use one to add content around the document, apply your own styles, or move parts of the document around.

Amatl ships three layouts. Pass another one with the `--html-layout` flag.

## Built-in layouts

All three layouts load `amatl://components.css`, which styles the `columns`, `grid`, `card` and `callout` container directives. See ["Container directives"](../directives/README.md#container-directives).

### `amatl://document.html`

A print-friendly layout for reports and other text-heavy documents, in A4 by default. It is the default layout for the `html` and `pdf` commands.

[View a sample document (PDF)](../../examples/document/document.pdf)

### `amatl://presentation.html`

Turns the document into slides. A horizontal rule (`---`) starts a new slide when it sits at the top level of the document. Anywhere else, for example in a container directive, a blockquote or a list item, `---` stays a visible rule in the current slide.

[View a sample presentation (PDF)](../../examples/presentation/presentation.pdf)

### `amatl://website.html`

A responsive layout for publishing Markdown as a web page.

[Visit the example website](https://bornholm.github.io/amatl/)

### Page title

The three layouts take the page title from the first of these that is set:

1. the `title` layout variable, given with `--html-layout-vars`;
2. the `title` key of the document front matter;
3. the text of the first `h1` of the document.

## Using a custom layout

Give the path or URL of your layout to the `--html-layout` flag:

```sh
amatl render html -o output.html --html-layout file://my-layout.html my-doc.md
```

`resolve` looks up relative paths from the directory of the layout, not from the directory of the document.

### Template data

A layout is a Go [`html/template`](https://pkg.go.dev/html/template). It receives:

| Field      | Content                                                                 |
| ---------- | ----------------------------------------------------------------------- |
| `.Body`    | the rendered document, as HTML                                          |
| `.Meta`    | the document front matter                                               |
| `.Vars`    | the layout variables given with `--html-layout-vars`                    |
| `.Context` | the rendering context, to pass to `resolve`                             |

### Functions

Every [Sprig](https://masterminds.github.io/sprig/) function is available, plus the following ones. The `html*` functions take an HTML string and a CSS selector.

| Function                                  | Returns                                                                                   |
| ----------------------------------------- | ----------------------------------------------------------------------------------------- |
| `htmlQueryFirst .Body "h1"`               | the first matching element, as HTML, or an empty string                                   |
| `htmlQueryAll .Body "h2"`                 | every matching element, as a list of HTML strings                                         |
| `htmlTextContent .Body "h1:first-of-type"` | the text of the matching elements, without tags. When several elements match, their texts run together with no space between them |
| `htmlRemove .Body "h1:first-of-type"`     | the HTML without the matching elements                                                    |
| `htmlAddAttr .Body "table" "class" "wide"` | the HTML with the value appended after a space to the attribute of each matching element, or the attribute created. Use it for attributes that hold a list, such as `class` or `rel`. On `id`, it gives `id="data wide"` |
| `htmlSplit .Body "hr"`                    | the top-level content cut at each matching element, as a list of HTML strings             |
| `resolve .Context "./logo.svg" "image/svg+xml"` | the resource as a `data:` URL, with the given MIME type or a detected one          |

`htmlQueryFirst`, `htmlRemove` and `htmlAddAttr` return a fragment that you can insert anywhere in the layout. An element the HTML parser moves to the head, such as a `<style>` written at the top of the document, stays in the result.

This layout puts the first heading in a banner and the rest of the document below it:

```html
<header class="banner">{{"{{"}} htmlQueryFirst .Body "h1:first-of-type" {{"}}"}}</header>
<main>{{"{{"}} htmlRemove .Body "h1:first-of-type" {{"}}"}}</main>
```
