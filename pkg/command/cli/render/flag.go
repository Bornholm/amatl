package render

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Bornholm/amatl/pkg/html"
	"github.com/Bornholm/amatl/pkg/html/layout/embed"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

const (
	paramToc            = "toc"
	paramHTMLLayout     = "html-layout"
	paramHTMLLayoutVars = "html-layout-vars"
)

var (
	flagToc = &cli.BoolFlag{
		Name:  paramToc,
		Value: false,
	}
	flagHTMLLayout = &cli.StringFlag{
		Name:  paramHTMLLayout,
		Value: html.DefaultLayoutURL,
		Usage: fmt.Sprintf("html layout to use, embedded: %v", embed.Available()),
	}
	flagHTMLLayoutVars = &cli.StringFlag{
		Name:  paramHTMLLayoutVars,
		Value: "{}",
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
