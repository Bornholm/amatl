# Templating

> **Available for:** `Markdown`, `HTML`, `PDF`

You can leverage [Go templating](https://pkg.go.dev/text/template) within your Markdown documents to inject dynamic content during generation.

To activate this functionality, simply utilize the `--vars` option along with a JSON object containing your desired values when executing the `render` command.

For instance, consider the following command:

```
amatl render pdf --vars '{"foo":"bar"}' my-doc.md
```

In `my-doc.md`, you can incorporate the injected value as follows:

```markdown
# My document

Here my value will be replaced: {{ .Vars.foo }}
```

For added convenience, amatl provides access to the [sprig](https://masterminds.github.io/sprig/) function library within templates.