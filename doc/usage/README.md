# Usage

Amatl's `render` command allows you to convert Markdown files into various output formats. Below are common use cases.

## 🖥️ Generate an HTML file

```sh
amatl render html -o output.html your-file.md
```

This will convert `your-file.md` into a standalone `output.html` file.

## 🖨️ Generate a PDF file

> Note: PDF generation requires Chrome or Chromium to be installed on your system.

```sh
amatl render pdf -o output.pdf your-file.md
```

This creates a `output.pdf` file from the specified Markdown input.

## 📝 Generate a Markdown file (processed)

> Useful for combining multiple files using the `include{}` directive or for generating a table of contents using `toc{}`.

```sh
amatl render markdown -o output.md your-file.md
```

This produces a processed Markdown file with all directives resolved.
