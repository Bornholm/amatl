---
title: Amatl - Markdown processing utility
menu:
  - label: Getting started
    items:
      - label: Installation
        href: "#installation"
        items:
          - label: On Linux
            href: "#-on-linux"
          - label: On Windows
            href: "#-on-windows"
          - label: On MacOS
            href: "#-on-macos"
      - label: Usage
        href: "#usage"
  - label: Features
    items:
      - label: "URL resolving"
        href: "#url-resolving"
      - label: "Variables and templating"
        href: "#variables-and-templating"
      - label: "Directives"
        href: "#directives"
        items:
          - label: "include{}"
            href: "#includeurlurl-fromheadingsheadinglevel-shiftheadingslevelshift"
          - label: "toc{}"
            href: "#tocminlevelminlevel-maxlevelmaxlevel"
          - label: "attrs{}"
            href: "#attrsattributes"
      - label: "Layouts"
        href: "#layouts"
  - label: MCP Server
    items:
      - label: "Starting the server"
        href: "#-starting-the-server"
      - label: "Configuring your client"
        href: "#-configuring-your-mcp-client"
      - label: "Available tools"
        href: "#-available-tools"
  - label: How to
    items:
      - label: "Write your own layout"
        href: "#write-your-own-layout"
      - label: "Share your configuration"
        href: "#share-your-configuration"
      - label: "Use in CI"
        href: "#use-in-ci"
  - label: Misc
    items:
      - label: Attributions
        href: "#-attributions"
---

<style>
.logo {
  text-align: center;
}

.logo > img {
  width: 150px;
}
</style>

:attrs{class="logo"}

![](../resources/logo.svg)

## What's Amatl ?

Amatl is a simple command-line utility that can help you to transform your [CommonMark](https://commonmark.org/) (also known as [Markdown](https://fr.wikipedia.org/wiki/Markdown)) files into full-fledged HTML/PDF documents.

For example, this simple website is [generated with Amatl itself](./index.md).

### Features

- Create document from local or remote resources via [URL resolving](#url-resolving);
- Integrate [MermaidJS](https://mermaid.js.org/) diagrams et code blocks with syntax highlighting;
- Use [custom directives](#directives) to include others documents or generate tables of content;
- Use [pre-defined or custom layouts](#layouts) to transform your content into presentations, report, etc
- Use [Go templating](#variables-and-templating) to inject dynamic data into your document;
- Rewrite you relative links and embed external resources;
- Expose your Markdown workspace to AI assistants via the built-in [MCP server](#mcp-server).

> **Why the name `amatl` ?**
>
> Amate (Spanish: amate `[aˈmate]` from Nahuatl languages: āmatl `[ˈaːmat͡ɬ]` is a type of bark paper that has been manufactured in Mexico since the precontact times. It was used primarily to create codices.
>
> Source: [Wikipédia](https://en.wikipedia.org/wiki/Amate)

Amatl is a [free software project](https://github.com/Bornholm/amatl) published under the [MIT licence](../../LICENCE).

## Getting started

### Installation

:include{url="../../doc/install/linux.md", shiftHeadings="3"}

:include{url="../../doc/install/windows.md", shiftHeadings="3"}

:include{url="../../doc/install/macos.md", shiftHeadings="3"}

---

:include{url="../../doc/usage/README.md", shiftHeadings="2"}

---

:include{url="../../doc/url-resolving/README.md", shiftHeadings="1"}

---

:include{url="../../doc/templating/README.md", shiftHeadings="1"}

---

:include{url="../../doc/directives/README.md", shiftHeadings="1"}

---

:include{url="../../doc/layouts/README.md", shiftHeadings="1"}

---

:include{url="../../doc/mcp/README.md", shiftHeadings="1"}

---

## How to

### Write your own layout

> Coming soon...

### Share your configuration

> Coming soon...

### Use in CI

> Coming soon...

## 🙏 Attributions

This project wouldn’t be possible without the incredible work of the open-source community. Special thanks to:

- [`github.com/yuin/goldmark` and its satellites libraries](https://github.com/yuin/goldmark)
- [`github.com/jgthms/bulma`](https://github.com/jgthms/bulma)

**Thank you for your amazing work !**
