package render

import (
	"io"
	"log/slog"

	"github.com/Bornholm/amatl/pkg/log"
	"github.com/Bornholm/amatl/pkg/markdown/directive/attrs"
	"github.com/Bornholm/amatl/pkg/markdown/directive/toc"
	"github.com/Bornholm/amatl/pkg/pipeline"
	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func HTML() *cli.Command {
	flags := withHTMLFlags()
	return &cli.Command{
		Name:   "html",
		Flags:  flags,
		Before: altsrc.InitInputSourceWithContext(flags, NewResolverSourceFromFlagFunc("config")),
		Action: func(ctx *cli.Context) error {
			sourcePath, source, err := getMarkdownSource(ctx)
			if err != nil {
				return errors.WithStack(err)
			}

			layoutVars, err := getVars(ctx, paramHTMLLayoutVars)
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

			baseDir := sourcePath.Dir()
			sourcePath = sourcePath.Base()

			pipelineCtx := log.WithAttrs(ctx.Context, slog.Any("source", sourcePath.String()))
			pipelineCtx = resolver.WithWorkDir(pipelineCtx, baseDir)

			transformer := pipeline.Pipeline(
				// Preprocess the markdown entrypoint
				// document to include potential directives
				MarkdownMiddleware(
					WithSourcePath(sourcePath),
					WithLinkReplacements(linkReplacements),
					WithIgnoredDirectives(toc.Type, attrs.Type),
				),
				TemplateMiddleware(
					WithVars(vars),
					WithDelimiters(leftDelimiter, rightDelimiter),
				),
				// Render the consolidated document
				// as HTML
				HTMLMiddleware(
					WithMarkdownTransformerOptions(
						WithSourcePath(sourcePath),
						WithLinkReplacements(linkReplacements),
					),
					WithLayoutURL(getHTMLLayout(ctx)),
					WithLayoutVars(layoutVars),
				),
			)

			payload := pipeline.NewPayload(source)

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
