# Layouts

> **Available for:** `HTML`, `PDF`

A **layout** is an HTML template that can be used to decorate your rendered Markdown content.

Layouts allow you to:

- Insert content before or after the main Markdown content.
- Apply custom styling or structural transformations to match your organization's branding or formatting needs.
- Use layout-specific variables to further enhance the output.

Amatl provides **three built-in layouts**, and you can also supply your own custom layout using the `--html-layout` flag.

## ✨ Built-in layouts

### `amatl://document.html`

A clean, print-friendly layout designed for general documents (e.g. reports or papers), formatted in A4 by default.

- **Default layout** for `HTML` and `PDF` outputs.
- Best suited for structured, text-heavy documents.

**Example:**  
[📄 View sample document (PDF)](../../examples/document/document.pdf)

**Variables:**

- `title`: Sets the HTML document title.

### `amatl://presentation.html`

A layout for creating slide-style presentations from your Markdown content.

**Example:**  
[📊 View sample presentation (PDF)](../../examples/presentation/presentation.pdf)

**Variables:**

- `title`: Sets the HTML document title.

---

### `amatl://website.html`

A simple, responsive layout suitable for rendering Markdown as a web page.

**Example:**  
[🌐 Visit example website](https://bornholm.github.io/amatl/)

**Variables:**

- `title`: Sets the HTML document title.

---

## 🛠️ Using a custom layout

To use a custom layout, provide the path or URL with the `--html-layout` flag:

```sh
amatl render markdown -o output.html --html-layout file://my-layout.html my-doc.md
```
