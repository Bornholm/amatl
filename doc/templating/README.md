# Templating

> **Available for:** `Markdown`, `HTML`, `PDF`

You can use [Go templating](https://pkg.go.dev/text/template) in your Markdown documents to inject dynamic content at generation time.

To enable this feature, simply pass the `--vars` with the JSON object containing your values to the `render` command.

For example, if i have the following command:

```
amatl render pdf --vars '{"foo":"bar"}' my-doc.md
```

In `my-doc.md` i can write the following to use the injected value:

```markdown
# My document

Here my value will be replaced: {{ .Vars.foo }}
```

For convenience, `amatl` exposes the [`sprig`](https://masterminds.github.io/sprig/) function library to templates.