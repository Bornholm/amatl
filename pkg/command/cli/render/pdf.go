package render

import (
	"io"
	"log/slog"

	"github.com/Bornholm/amatl/pkg/log"
	"github.com/Bornholm/amatl/pkg/pipeline"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func PDF() *cli.Command {
	flags := withPDFFlags()
	return &cli.Command{
		Name:   "pdf",
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

			linkReplacements, err := getLinkReplacements(ctx)
			if err != nil {
				return errors.WithStack(err)
			}

			layoutVars, err := getVars(ctx, paramHTMLLayoutVars)
			if err != nil {
				return errors.WithStack(err)
			}

			marginTop, marginRight, marginBottom, marginLeft := getPDFMargin(ctx)
			scale := getPDFScale(ctx)
			timeout := getPDFTimeout(ctx)
			background := getPDFBackground(ctx)
			execPath := getPDFExecPath(ctx)

			transformer := pipeline.Pipeline(
				// Preprocess the markdown entrypoint
				// document to include potential directives
				MarkdownMiddleware(
					WithSourceURL(sourceURL),
					WithLinkReplacements(linkReplacements),
				),
				TemplateMiddleware(
					WithVars(vars),
				),
				// Render the consolidated document
				// as HTML
				HTMLMiddleware(
					WithMarkdownTransformerOptions(
						WithSourceURL(sourceURL),
						WithLinkReplacements(linkReplacements),
					),
					WithLayoutURL(getHTMLLayout(ctx)),
					WithLayoutVars(layoutVars),
				),
				// Render generated HTML to PDF with Chromium
				PDFMiddleware(
					WithMarginTop(marginTop),
					WithMarginRight(marginRight),
					WithMarginBottom(marginBottom),
					WithMarginLeft(marginLeft),
					WithScale(scale),
					WithTimeout(timeout),
					WithBackground(background),
					WithExecPath(execPath),
				),
			)

			payload := pipeline.NewPayload(source)

			pipelineCtx := log.WithAttrs(ctx.Context, slog.Any("source", sourceURL.String()))

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
