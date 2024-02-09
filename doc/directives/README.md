# Directives

> **Available for:** `Markdown`, `HTML`, `PDF` 

`amatl` extends the CommonMark syntax with what we call "directives".

A directive takes this form:

```
:directive{attr="value1" attr="value2" ...}
```

Each directive triggers a specific behavior based on its type.

## Available directives

### `:include{path="<path>"}`

Include another Markdown file in your document.

