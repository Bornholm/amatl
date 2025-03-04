package render

import (
	"io"

	"github.com/Bornholm/amatl/pkg/pipeline"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func Markdown() *cli.Command {
	flags := withCommonFlags()
	return &cli.Command{
		Name:   "markdown",
		Flags:  flags,
		Before: altsrc.InitInputSourceWithContext(flags, NewResolverSourceFromFlagFunc("config")),
		Action: func(ctx *cli.Context) error {
			sourceURL, source, err := getMarkdownSource(ctx)
			if err != nil {
				return errors.WithStack(err)
			}

			vars, err := getVars(ctx, paramVars)
			if err != nil {
				return errors.WithStack(err)
			}

			transformer := pipeline.Pipeline(
				MarkdownMiddleware(
					WithSourceURL(sourceURL),
				),
				TemplateMiddleware(
					WithVars(vars),
				),
			)

			payload := pipeline.NewPayload(ctx.Context, source)

			if err := transformer.Transform(payload); err != nil {
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

			if _, err := io.Copy(output, payload.Buffer()); err != nil {
				return errors.WithStack(err)
			}

			return nil
		},
	}
}
