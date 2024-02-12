package render

import (
	"bytes"
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

			vars, err := getVars(ctx)
			if err != nil {
				return errors.WithStack(err)
			}

			layoutVars, err := getHTMLLayoutVars(ctx)
			if err != nil {
				return errors.WithStack(err)
			}

			marginTop, marginRight, marginBottom, marginLeft := getPDFMargin(ctx)
			scale := getPDFScale(ctx)

			pipeline := pipeline.New(
				// Preprocess the markdown entrypoint
				// document to include potential directives
				MarkdownTransformer(
					WithSourceURL(sourceURL),
					WithToc(false),
				),
				ToggleableTransformer(
					TemplateTransformer(
						WithVars(vars),
					),
					hasVars(ctx),
				),
				// Render the consolidated document
				// as HTML
				HTMLTransformer(
					WithMarkdownTransformerOptions(
						WithSourceURL(sourceURL),
						WithToc(isTocEnabled(ctx)),
					),
					WithLayoutURL(getHTMLLayout(ctx)),
					WithLayoutVars(layoutVars),
				),
				// Render generated HTML to PDF with Chromium
				PDFTransformer(
					WithMarginTop(marginTop),
					WithMarginRight(marginRight),
					WithMarginBottom(marginBottom),
					WithMarginLeft(marginLeft),
					WithScale(scale),
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
