package render

import (
	"bytes"
	"io"

	"github.com/Bornholm/amatl/pkg/pipeline"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func Markdown() *cli.Command {
	return &cli.Command{
		Name:  "markdown",
		Flags: withCommonFlags(),
		Action: func(ctx *cli.Context) error {
			_, dirname, source, err := getMarkdownSource(ctx)
			if err != nil {
				return errors.WithStack(err)
			}

			vars, err := getVars(ctx)
			if err != nil {
				return errors.WithStack(err)
			}

			pipeline := pipeline.New(
				ToggleableTransformer(
					TemplateTransformer(
						WithVars(vars),
					),
					hasVars(ctx),
				),
				MarkdownTransformer(
					WithBaseDir(dirname),
					WithToc(isTocEnabled(ctx)),
				),
			)

			result, err := pipeline.Transform(ctx.Context, source)
			if err != nil {
				return errors.WithStack(err)
			}

			output, err := getOutput(ctx)
			if err != nil {
				return errors.WithStack(err)
			}

			defer func() {
				if err := output.Close(); err != nil {
					panic(errors.WithStack(err))
				}
			}()

			if _, err := io.Copy(output, bytes.NewBuffer(result)); err != nil {
				return errors.WithStack(err)
			}

			return nil
		},
	}
}
