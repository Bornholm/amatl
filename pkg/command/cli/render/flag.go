package render

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Bornholm/amatl/pkg/html/layout"
	"github.com/Bornholm/amatl/pkg/html/layout/resolver/amatl"
	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/Bornholm/amatl/pkg/transform"
	"gopkg.in/yaml.v3"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	// Register resolver schemes

	_ "github.com/Bornholm/amatl/pkg/resolver/file"
	_ "github.com/Bornholm/amatl/pkg/resolver/http"
	_ "github.com/Bornholm/amatl/pkg/resolver/stdin"
)

const (
	paramOutput                 = "output"
	paramTemplateVars           = "vars"
	paramTemplateLeftDelimiter  = "template-left-delimiter"
	paramTemplateRightDelimiter = "template-right-delimiter"
	paramLinkReplacements       = "link-replacements"
	paramHTMLLayout             = "html-layout"
	paramHTMLLayoutVars         = "html-layout-vars"
	paramPDFMarginTop           = "pdf-margin-top"
	paramPDFMarginLeft          = "pdf-margin-left"
	paramPDFMarginRight         = "pdf-margin-right"
	paramPDFMarginBottom        = "pdf-margin-bottom"
	paramPDFScale               = "pdf-scale"
	paramPDFTimeout             = "pdf-timeout"
	paramPDFBackground          = "pdf-background"
	paramPDFExecPath            = "pdf-exec-path"
	paramPDFDisplayHeaderFooter = "pdf-display-header-footer"
	paramPDFHeaderTemplate      = "pdf-header-template"
	paramPDFFooterTemplate      = "pdf-footer-template"
	paramPDFNoSandbox           = "pdf-no-sandbox"
)

var (
	flagOutput = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    paramOutput,
		Aliases: []string{"o"},
		Value:   "-",
		Usage:   "output generated content to given file, '-' to write to stdout",
	})
	flagTemplateVars = altsrc.NewStringFlag(&cli.StringFlag{
		Name:  paramTemplateVars,
		Value: "",
		Usage: "enable templating and use url resource as json injected data",
	})
	flagTemplateLeftDelimiter = altsrc.NewStringFlag(&cli.StringFlag{
		Name:  paramTemplateLeftDelimiter,
		Value: "",
		Usage: "override default templating left delimiter",
	})
	flagTemplateRightDelimiter = altsrc.NewStringFlag(&cli.StringFlag{
		Name:  paramTemplateRightDelimiter,
		Value: "",
		Usage: "override default templating right delimiter",
	})
	flagHTMLLayout = altsrc.NewStringFlag(&cli.StringFlag{
		Name:  paramHTMLLayout,
		Value: layout.DefaultRawURL,
		Usage: fmt.Sprintf("html layout to use, available by default: %v", amatl.Available()),
	})
	flagHTMLLayoutVars = altsrc.NewStringFlag(&cli.StringFlag{
		Name:  paramHTMLLayoutVars,
		Usage: "enable layout templating and use url resource as json injected data",
		Value: "",
	})
	flagLinkReplacements = altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
		Name:  paramLinkReplacements,
		Usage: "replace the given link prefix by the string provided, expected format <prefix>::<replacement>",
		Value: cli.NewStringSlice(),
	})
	flagPDFMarginTop = altsrc.NewFloat64Flag(&cli.Float64Flag{
		Name:  paramPDFMarginTop,
		Value: DefaultPDFMargin,
		Usage: "pdf top margin in centimeters",
	})
	flagPDFMarginRight = altsrc.NewFloat64Flag(&cli.Float64Flag{
		Name:  paramPDFMarginRight,
		Value: DefaultPDFMargin,
		Usage: "pdf right margin in centimeters",
	})
	flagPDFMarginLeft = altsrc.NewFloat64Flag(&cli.Float64Flag{
		Name:  paramPDFMarginLeft,
		Value: DefaultPDFMargin,
		Usage: "pdf left margin in centimeters",
	})
	flagPDFMarginBottom = altsrc.NewFloat64Flag(&cli.Float64Flag{
		Name:  paramPDFMarginBottom,
		Value: DefaultPDFMargin,
		Usage: "pdf bottom margin in centimeters",
	})
	flagPDFScale = altsrc.NewFloat64Flag(&cli.Float64Flag{
		Name:  paramPDFScale,
		Value: DefaultPDFScale,
		Usage: "pdf print scale",
	})
	flagPDFTimeout = altsrc.NewDurationFlag(&cli.DurationFlag{
		Name:  paramPDFTimeout,
		Value: DefaultPDFTimeout,
		Usage: "pdf generation timeout",
	})
	flagPDFBackground = altsrc.NewBoolFlag(&cli.BoolFlag{
		Name:  paramPDFBackground,
		Value: DefaultPDFBackground,
		Usage: "pdf print background",
	})
	flagPDFExecPath = altsrc.NewStringFlag(&cli.StringFlag{
		Name:  paramPDFExecPath,
		Value: DefaultPDFExecPath,
		Usage: "pdf chromium executable path",
	})
	flagPDFDisplayHeaderFooter = altsrc.NewBoolFlag(&cli.BoolFlag{
		Name:  paramPDFDisplayHeaderFooter,
		Usage: "pdf display header and footer",
		Value: false,
	})
	flagPDFFooterTemplate = altsrc.NewStringFlag(&cli.StringFlag{
		Name:  paramPDFFooterTemplate,
		Usage: "pdf footer template",
		Value: DefaultPDFFooterTemplate,
	})
	flagPDFFHeaderTemplate = altsrc.NewStringFlag(&cli.StringFlag{
		Name:  paramPDFHeaderTemplate,
		Usage: "pdf header template",
		Value: DefaultPDFHeaderTemplate,
	})
	flagPDFNoSandbox = altsrc.NewStringFlag(&cli.StringFlag{
		Name:  paramPDFNoSandbox,
		Usage: "disable chrome sandboxing",
		Value: DefaultPDFHeaderTemplate,
	})
)

func getVars(ctx *cli.Context, param string) (map[string]any, error) {
	rawUrl := ctx.String(param)

	if rawUrl == "" {
		return map[string]any{}, nil
	}

	path := resolver.Path(rawUrl)

	reader, err := resolver.Resolve(ctx.Context, path.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer func() {
		if err := reader.Close(); err != nil {
			panic(errors.WithStack(err))
		}
	}()

	transformed := transform.NewNewlineReader(reader)

	source, err := io.ReadAll(transformed)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var vars map[string]any
	if err := json.Unmarshal([]byte(source), &vars); err != nil {
		return nil, errors.Wrap(err, "could not parse vars")
	}

	return vars, nil
}

func getTemplateDelimiters(ctx *cli.Context) (string, string) {
	left := ctx.String(paramTemplateLeftDelimiter)
	right := ctx.String(paramTemplateRightDelimiter)
	return left, right
}

func getHTMLLayout(ctx *cli.Context) string {
	return ctx.String(paramHTMLLayout)
}

func withCommonFlags(flags ...cli.Flag) []cli.Flag {
	return append([]cli.Flag{
		flagTemplateVars,
		flagTemplateLeftDelimiter,
		flagTemplateRightDelimiter,
		flagOutput,
		flagLinkReplacements,
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
		flagPDFTimeout,
		flagPDFBackground,
		flagPDFExecPath,
		flagPDFDisplayHeaderFooter,
		flagPDFFHeaderTemplate,
		flagPDFFooterTemplate,
		flagPDFNoSandbox,
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

func getLinkReplacements(ctx *cli.Context) (map[string]string, error) {
	rawLinkReplacements := ctx.StringSlice(paramLinkReplacements)

	linkReplacements := make(map[string]string)
	for _, r := range rawLinkReplacements {
		parts := strings.SplitN(r, "::", 2)
		if len(parts) != 2 {
			return nil, errors.Errorf("invalid link replacement format '%s'", r)
		}

		linkReplacements[parts[0]] = parts[1]
	}

	return linkReplacements, nil
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

	path := resolver.Path(filename)

	reader, err := resolver.Resolve(ctx.Context, path.String())
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	defer func() {
		if err := reader.Close(); err != nil {
			panic(errors.WithStack(err))
		}
	}()

	transformed := transform.NewNewlineReader(reader)

	source, err := io.ReadAll(transformed)
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

func getPDFTimeout(ctx *cli.Context) time.Duration {
	return ctx.Duration(paramPDFTimeout)
}

func getPDFBackground(ctx *cli.Context) bool {
	return ctx.Bool(paramPDFBackground)
}

func getPDFExecPath(ctx *cli.Context) string {
	return ctx.String(paramPDFExecPath)
}

func getPDFHeaderFooter(ctx *cli.Context) (bool, string, string) {
	return ctx.Bool(paramPDFDisplayHeaderFooter),
		ctx.String(paramPDFHeaderTemplate),
		ctx.String(paramPDFFooterTemplate)
}

func getPDFNoSandbox(ctx *cli.Context) bool {
	return ctx.Bool(paramPDFNoSandbox)
}

func NewResolverSourceFromFlagFunc(flag string) func(cCtx *cli.Context) (altsrc.InputSourceContext, error) {
	return func(cCtx *cli.Context) (altsrc.InputSourceContext, error) {
		if urlStr := cCtx.String(flag); urlStr != "" {
			return NewResolvedInputSource(cCtx.Context, urlStr)
		}

		return altsrc.NewMapInputSource("", map[interface{}]interface{}{}), nil
	}
}

func NewResolvedInputSource(ctx context.Context, urlStr string) (altsrc.InputSourceContext, error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse url '%s'", urlStr)
	}

	path := resolver.Path(urlStr)

	reader, err := resolver.Resolve(ctx, path.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer func() {
		if err := reader.Close(); err != nil {
			panic(errors.WithStack(err))
		}
	}()

	transformed := transform.NewNewlineReader(reader)

	data, err := io.ReadAll(transformed)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ext := filepath.Ext(urlStr)
	switch ext {
	case ".json":
		fallthrough
	case ".yaml":
		fallthrough
	case ".yml":
		var values map[any]any

		if err := yaml.Unmarshal(data, &values); err != nil {
			return nil, errors.WithStack(err)
		}

		values, err = rewriteRelativeURL(url, values)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return altsrc.NewMapInputSource(urlStr, values), nil

	default:
		return nil, errors.Errorf("no parser associated with '%s' file extension", ext)
	}

}

func rewriteRelativeURL(fromURL *url.URL, values map[any]any) (map[any]any, error) {
	sourcePath := resolver.Path(fromURL.String())
	sourceDir := sourcePath.Dir()

	for key, rawValue := range values {
		value, ok := rawValue.(string)
		if !ok {
			continue
		}

		switch {
		case isURL(value):
			continue

		case isPath(value):
			if filepath.IsAbs(value) {
				continue
			}

			abs := sourceDir.JoinPath(value)
			values[key] = abs.String()
			continue
		}

	}

	return values, nil
}

var filepathRegExp = regexp.MustCompile(`^(?i)(?:\/[^\/]+)+\/?[^\s]+(?:\.[^\s]+)+|[^\s]+(?:\.[^\s]+)+$`)

func isPath(str string) bool {
	return filepathRegExp.MatchString(str)
}

func isURL(str string) bool {
	_, err := url.ParseRequestURI(str)
	return err == nil
}
