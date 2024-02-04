package render

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Bornholm/amatl/pkg/html"
	"github.com/Bornholm/amatl/pkg/html/layout/base"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

const (
	paramToc             = "toc"
	paramHTMLLayout      = "html-layout"
	paramHTMLLayoutVars  = "html-layout-vars"
	paramPDFMarginTop    = "pdf-margin-top"
	paramPDFMarginLeft   = "pdf-margin-left"
	paramPDFMarginRight  = "pdf-margin-right"
	paramPDFMarginBottom = "pdf-margin-bottom"
	paramPDFScale        = "pdf-scale"
)

var (
	flagToc = &cli.BoolFlag{
		Name:  paramToc,
		Value: false,
	}
	flagHTMLLayout = &cli.StringFlag{
		Name:  paramHTMLLayout,
		Value: html.DefaultLayoutURL,
		Usage: fmt.Sprintf("html layout to use, embedded: %v", base.Available()),
	}
	flagHTMLLayoutVars = &cli.StringFlag{
		Name:  paramHTMLLayoutVars,
		Value: "{}",
	}
	flagPDFMarginTop = &cli.Float64Flag{
		Name:  paramPDFMarginTop,
		Value: DefaultPDFMargin,
		Usage: "pdf top margin in centimeters",
	}
	flagPDFMarginRight = &cli.Float64Flag{
		Name:  paramPDFMarginRight,
		Value: DefaultPDFMargin,
		Usage: "pdf right margin in centimeters",
	}
	flagPDFMarginLeft = &cli.Float64Flag{
		Name:  paramPDFMarginLeft,
		Value: DefaultPDFMargin,
		Usage: "pdf left margin in centimeters",
	}
	flagPDFMarginBottom = &cli.Float64Flag{
		Name:  paramPDFMarginBottom,
		Value: DefaultPDFMargin,
		Usage: "pdf bottom margin in centimeters",
	}
	flagPDFScale = &cli.Float64Flag{
		Name:  paramPDFScale,
		Value: DefaultPDFScale,
		Usage: "pdf print scale",
	}
)

func isTocEnabled(ctx *cli.Context) bool {
	return ctx.Bool(paramToc)
}

func getHTMLLayout(ctx *cli.Context) string {
	return ctx.String(paramHTMLLayout)
}

func getHTMLLayoutVars(ctx *cli.Context) (map[string]any, error) {
	rawVars := ctx.String(paramHTMLLayoutVars)

	var vars map[string]any
	if err := json.Unmarshal([]byte(rawVars), &vars); err != nil {
		return nil, errors.Wrap(err, "could not parse html layout vars")
	}

	return vars, nil
}

func withCommonFlags(flags ...cli.Flag) []cli.Flag {
	return append([]cli.Flag{
		flagToc,
	}, flags...)
}

func withHTMLFlags(flags ...cli.Flag) []cli.Flag {
	flags = append(flags,
		flagHTMLLayout,
		flagHTMLLayoutVars,
	)

	return withCommonFlags(flags...)
}

func withPDFFlags(flags ...cli.Flag) []cli.Flag {
	flags = append(flags,
		flagPDFMarginTop,
		flagPDFMarginLeft,
		flagPDFMarginRight,
		flagPDFMarginBottom,
		flagPDFScale,
	)

	return withHTMLFlags(flags...)
}

func getMarkdownSource(ctx *cli.Context) (string, string, []byte, error) {
	filename := ctx.Args().First()
	dirname, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		return "", "", nil, errors.WithStack(err)
	}

	source, err := os.ReadFile(filename)
	if err != nil {
		return "", "", nil, errors.WithStack(err)
	}

	return filename, dirname, source, nil
}

func getPDFScale(ctx *cli.Context) float64 {
	return ctx.Float64(paramPDFScale)
}

func getPDFMargin(ctx *cli.Context) (top float64, right float64, bottom float64, left float64) {
	return ctx.Float64(paramPDFMarginTop),
		ctx.Float64(paramPDFMarginRight),
		ctx.Float64(paramPDFMarginBottom),
		ctx.Float64(paramPDFMarginLeft)
}
