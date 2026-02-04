package render

import (
	"io"
	"log/slog"

	"github.com/Bornholm/amatl/pkg/log"
	"github.com/Bornholm/amatl/pkg/pipeline"
	"github.com/Bornholm/amatl/pkg/resolver"
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
			sourcePath, source, err := getMarkdownSource(ctx)
			if err != nil {
				return errors.WithStack(err)
			}

			vars, err := getVars(ctx, paramTemplateVars)
			if err != nil {
				return errors.WithStack(err)
			}

			leftDelimiter, rightDelimiter := getTemplateDelimiters(ctx)

			linkReplacements, err := getLinkReplacements(ctx)
			if err != nil {
				return errors.WithStack(err)
			}

			transformer := pipeline.Pipeline(
				MarkdownMiddleware(
					WithSourcePath(sourcePath),
					WithLinkReplacements(linkReplacements),
				),
				TemplateMiddleware(
					WithVars(vars),
					WithDelimiters(leftDelimiter, rightDelimiter),
				),
			)

			payload := pipeline.NewPayload(source)

			baseDir := sourcePath.Dir()
			sourcePath = sourcePath.Base()

			pipelineCtx := log.WithAttrs(ctx.Context, slog.Any("source", sourcePath.String()))
			pipelineCtx = resolver.WithWorkDir(pipelineCtx, baseDir)

			if err := transformer.Transform(pipelineCtx, payload); err != nil {
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
