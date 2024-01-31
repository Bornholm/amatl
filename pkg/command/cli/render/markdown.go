package render

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/Bornholm/amatl/pkg/pipeline"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func Markdown() *cli.Command {
	return &cli.Command{
		Name:  "markdown",
		Flags: withCommonFlags(),
		Action: func(ctx *cli.Context) error {
			for _, filename := range ctx.Args().Slice() {
				baseDir, err := filepath.Abs(filepath.Dir(filename))
				if err != nil {
					return errors.WithStack(err)
				}

				source, err := os.ReadFile(filename)
				if err != nil {
					return errors.WithStack(err)
				}

				pipeline := pipeline.New(
					MarkdownTransformer(
						WithBaseDir(baseDir),
						WithToc(isTocEnabled(ctx)),
					),
				)

				result, err := pipeline.Transform(ctx.Context, source)
				if err != nil {
					return errors.WithStack(err)
				}

				if _, err := io.Copy(os.Stdout, bytes.NewBuffer(result)); err != nil {
					return errors.WithStack(err)
				}
			}

			return nil
		},
	}
}
