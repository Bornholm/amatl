# Variables and templating

> **Available for:** `Markdown`, `HTML`, `PDF`

You can leverage [Go templating](https://pkg.go.dev/text/template) within your Markdown documents to inject dynamic content during generation.

To activate this functionality, simply utilize the `--vars` option along with an URL of a JSON resource containing your desired values when executing the `render` command.

For instance, consider the following command:

```shell
echo '{"foo":"bar"}' | amatl render pdf --vars stdin://  my-doc.md
```

In `my-doc.md`, you can incorporate the injected value as follows:

<!--
In the following example, the expected template delimiters are `{‎‎{` and `}‎}` (without the space in between).
They are escaped here for the website https://bornholm.github.io/amatl/.
-->

```markdown
# My document

Here my value will be replaced: {{"{{"}} .Vars.foo {{"}}"}}
```

For added convenience, amatl provides access to the [sprig](https://masterminds.github.io/sprig/) function library within templates.
