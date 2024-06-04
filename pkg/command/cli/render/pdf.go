package render

import (
	"io"

	"github.com/Bornholm/amatl/pkg/pipeline"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func PDF() *cli.Command {
	return &cli.Command{
		Name:  "pdf",
		Flags: withPDFFlags(),
		Action: func(ctx *cli.Context) error {
			sourceURL, source, err := getMarkdownSource(ctx)
			if err != nil {
				return errors.WithStack(err)
			}

			vars, err := getVars(ctx, paramVars)
			if err != nil {
				return errors.WithStack(err)
			}

			layoutVars, err := getVars(ctx, paramHTMLLayoutVars)
			if err != nil {
				return errors.WithStack(err)
			}

			marginTop, marginRight, marginBottom, marginLeft := getPDFMargin(ctx)
			scale := getPDFScale(ctx)

			transformer := pipeline.Pipeline(
				// Preprocess the markdown entrypoint
				// document to include potential directives
				MarkdownMiddleware(
					WithSourceURL(sourceURL),
				),
				TemplateMiddleware(
					WithVars(vars),
				),
				// Render the consolidated document
				// as HTML
				HTMLMiddleware(
					WithMarkdownTransformerOptions(
						WithSourceURL(sourceURL),
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
