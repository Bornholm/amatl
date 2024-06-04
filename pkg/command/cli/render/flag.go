package render

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/Bornholm/amatl/pkg/html/layout"
	"github.com/Bornholm/amatl/pkg/html/layout/resolver/amatl"
	"github.com/Bornholm/amatl/pkg/resolver"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	// Register resolver schemes

	_ "github.com/Bornholm/amatl/pkg/resolver/file"
	_ "github.com/Bornholm/amatl/pkg/resolver/http"
	_ "github.com/Bornholm/amatl/pkg/resolver/stdin"
)

const (
	paramVars            = "vars"
	paramOutput          = "output"
	paramHTMLLayout      = "html-layout"
	paramHTMLLayoutVars  = "html-layout-vars"
	paramPDFMarginTop    = "pdf-margin-top"
	paramPDFMarginLeft   = "pdf-margin-left"
	paramPDFMarginRight  = "pdf-margin-right"
	paramPDFMarginBottom = "pdf-margin-bottom"
	paramPDFScale        = "pdf-scale"
)

var (
	flagOutput = &cli.StringFlag{
		Name:    paramOutput,
		Aliases: []string{"o"},
		Value:   "-",
		Usage:   "output generated content to given file, '-' to write to stdout",
	}
	flagVars = &cli.StringFlag{
		Name:  paramVars,
		Value: "",
		Usage: "enable templating and use url resource as json injected data",
	}
	flagHTMLLayout = &cli.StringFlag{
		Name:  paramHTMLLayout,
		Value: layout.DefaultRawURL,
		Usage: fmt.Sprintf("html layout to use, available by default: %v", amatl.Available()),
	}
	flagHTMLLayoutVars = &cli.StringFlag{
		Name:  paramHTMLLayoutVars,
		Usage: "enable layout templating and use url resource as json injected data",
		Value: "",
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

func getVars(ctx *cli.Context, param string) (map[string]any, error) {
	rawUrl := ctx.String(param)

	if rawUrl == "" {
		return map[string]any{}, nil
	}

	url, err := url.Parse(rawUrl)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	reader, err := resolver.Resolve(ctx.Context, url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer func() {
		if err := reader.Close(); err != nil {
			panic(errors.WithStack(err))
		}
	}()

	source, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var vars map[string]any
	if err := json.Unmarshal([]byte(source), &vars); err != nil {
		return nil, errors.Wrap(err, "could not parse vars")
	}

	return vars, nil
}

func getHTMLLayout(ctx *cli.Context) string {
	return ctx.String(paramHTMLLayout)
}

func withCommonFlags(flags ...cli.Flag) []cli.Flag {
	return append([]cli.Flag{
		flagVars,
		flagOutput,
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

func getOutput(ctx *cli.Context) (io.WriteCloser, error) {
	output := ctx.String(paramOutput)
	if output == "-" {
		return os.Stdout, nil
	}

	file, err := os.Create(output)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return file, nil
}

func getMarkdownSource(ctx *cli.Context) (*url.URL, []byte, error) {
	filename := ctx.Args().First()
	if filename == "" {
		return nil, nil, errors.New("you must provide the path or url to a markdown file")
	}

	url, err := url.Parse(filename)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	reader, err := resolver.Resolve(ctx.Context, url)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	defer func() {
		if err := reader.Close(); err != nil {
			panic(errors.WithStack(err))
		}
	}()

	source, err := io.ReadAll(reader)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return url, source, nil
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
