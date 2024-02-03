package render

import (
	"bytes"
	"io"
	"os"

	"github.com/Bornholm/amatl/pkg/pipeline"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func PDF() *cli.Command {
	return &cli.Command{
		Name:  "pdf",
		Flags: withPDFFlags(),
		Action: func(ctx *cli.Context) error {
			_, dirname, source, err := getMarkdownSource(ctx)
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
					WithBaseDir(dirname),
					WithToc(false),
				),
				// Render the consolidated document
				// as HTML
				HTMLTransformer(
					WithMarkdownTransformerOptions(
						WithBaseDir(dirname),
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

			if _, err := io.Copy(os.Stdout, bytes.NewBuffer(result)); err != nil {
				return errors.WithStack(err)
			}

			return nil
		},
	}
}
