package render

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"forge.cadoles.com/wpetit/amatl/pkg/markdown/directive"
	"forge.cadoles.com/wpetit/amatl/pkg/markdown/directive/include"
	"forge.cadoles.com/wpetit/amatl/pkg/markdown/renderer/markdown"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func Markdown() *cli.Command {
	return &cli.Command{
		Name:  "markdown",
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

				if _, err := io.Copy(os.Stdout, &buf); err != nil {
					return errors.WithStack(err)
				}
			}

			return nil
		},
	}
}

func newMarkdownRenderer() renderer.Renderer {
	render := markdown.NewRenderer()

	render.AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(
				directive.NewRenderer(
					directive.WithRenderer(
						include.Type,
						&include.NodeRenderer{
							Cache:    cache,
							Renderer: render,
						},
					),
				), 0,
			),
		),
		markdown.WithNodeRenderer(
			directive.KindDirective,
			directive.NewMarkdownNodeRenderer(
				directive.WithMarkdownDirectiveRenderer(
					include.Type,
					&include.MarkdownRenderer{
						Cache: cache,
					},
				),
			),
		),
	)

	return render
}
