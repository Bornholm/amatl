package render

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
)

func HTML() *cli.Command {
	return &cli.Command{
		Name:  "html",
		Flags: withCommonFlags(),
		Action: func(ctx *cli.Context) error {
			for _, filename := range ctx.Args().Slice() {
				basePath, err := filepath.Abs(filepath.Dir(filename))
				if err != nil {
					return errors.WithStack(err)
				}

				source, err := os.ReadFile(filename)
				if err != nil {
					return errors.WithStack(err)
				}

				reader := text.NewReader(source)

				parse := newParser(basePath, false)
				render := newMarkdownRenderer()

				document := parse.Parse(reader)

				var buf bytes.Buffer
				if err := render.Render(&buf, source, document); err != nil {
					return errors.WithStack(err)
				}

				source = buf.Bytes()
				reader = text.NewReader(source)

				withToc := isTocEnabled(ctx)

				parse = newParser(basePath, withToc)
				document = parse.Parse(reader)

				render = newHTMLRenderer()

				if err := render.Render(os.Stdout, source, document); err != nil {
					return errors.WithStack(err)
				}
			}

			return nil
		},
	}
}

func newHTMLRenderer() renderer.Renderer {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
		),
	)

	return markdown.Renderer()
}
