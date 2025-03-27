<p align="center">
  <img src="./misc/resources/logo.svg" style="height:150px" />
</p>

# `amatl`

Transform your [CommonMark](https://commonmark.org/) (also known as [Markdown](https://fr.wikipedia.org/wiki/Markdown)) files into full-fledged documents from the command line.

> **Why the name `amatl` ?**
>
> Amate (Spanish: amate `[aˈmate]` from Nahuatl languages: āmatl `[ˈaːmat͡ɬ]` is a type of bark paper that has been manufactured in Mexico since the precontact times. It was used primarily to create codices.
>
> Source: [Wikipédia](https://en.wikipedia.org/wiki/Amate)

## Features

- Create document from local or remote resources via URL resolving;
- Integrate [MermaidJS](https://mermaid.js.org/) diagrams et code blocks with syntax highlighting;
- Use [custom directives](./doc/directives/README.md) to include others documents or generate tables of content;
- Use [pre-defined or custom layouts](./doc/layouts/README.md) to transform your content into presentations, report, etc;
- Use [Go templating](./doc/templating/README.md) to inject dynamic data into your document.

## How to use

> **Caution ⚠**
>
> This project is in its early stages and subject to rapid evolution. Expect frequent changes and potential instability.

1. Download [the latest release](https://github.com/Bornholm/amatl/releases/latest) of `amatl`
2. In your terminal, start transforming your documents to HTML:

```shell
amatl render html my-doc.md
```

See [`./doc`](./doc/README.md) for more informations.

## Licence

[MIT](./LICENCE)
