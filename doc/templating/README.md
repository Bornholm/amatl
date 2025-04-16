# Variables and Templating

> **Available for:** `Markdown`, `HTML`, `PDF`

Amatl supports dynamic content injection using [Go templates](https://pkg.go.dev/text/template). This allows you to define reusable or customizable elements directly within your documents.

## 🔧 Injecting variables

Use the `--vars` option with a URL (e.g., a local file or `stdin://`) pointing to a JSON document containing your variable values. These will be accessible within the template under `.Vars`.

**Example:**

```sh
echo '{"foo": "bar"}' | amatl render pdf --vars stdin:// my-doc.md
```

In `my-doc.md`, you can access the variable like this:

<!-- Escaping the delimiters here for rendering on https://bornholm.github.io/amatl/ -->

```markdown
# My Document

Here my value will be replaced: {{"{{"}} .Vars.foo {{"}}"}}
```

Amatl also includes the [full Sprig function library](https://masterminds.github.io/sprig/) to enhance template expressions with string manipulation, logic functions, and more.

## 📄 Using metadata from YAML front matter

Amatl parses YAML front matter and exposes its content through the `.Meta` object in templates. This is particularly useful for static values defined in the document itself.

```markdown
---
foo: "bar"
---

# My Document

Here my value will be replaced: {{"{{"}} .Meta.foo {{"}}"}}
```
