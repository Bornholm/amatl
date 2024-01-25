package render

import "github.com/urfave/cli/v2"

const paramToc = "toc"

var (
	flagToc = &cli.BoolFlag{
		Name:  paramToc,
		Value: false,
	}
)

func isTocEnabled(ctx *cli.Context) bool {
	return ctx.Bool(paramToc)
}

func withCommonFlags(flags ...cli.Flag) []cli.Flag {
	return append([]cli.Flag{
		flagToc,
	}, flags...)
}
