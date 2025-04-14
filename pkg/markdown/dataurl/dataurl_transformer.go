package dataurl

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Bornholm/amatl/pkg/pipeline"
	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/pkg/errors"
	"github.com/vincent-petithory/dataurl"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type Transformer struct {
}

// Transform implements parser.ASTTransformer.
func (t *Transformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	err := ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		image, ok := n.(*ast.Image)
		if !ok {
			return ast.WalkContinue, nil
		}

		destination := string(image.Destination)

		resourceURL, err := url.Parse(destination)
		if err != nil {
			return ast.WalkStop, errors.WithStack(err)
		}

		ctx, err := pipeline.FromParserContext(pc)
		if err != nil {
			return ast.WalkStop, errors.WithStack(err)
		}

		resourceReader, err := resolver.Resolve(ctx, resourceURL)
		if err != nil {
			return ast.WalkStop, errors.Wrapf(err, "could not resolve resource '%s'", destination)
		}

		defer func() {
			if err := resourceReader.Close(); err != nil {
				panic(errors.Wrapf(err, "could not close resource '%s'", resourceURL))
			}
		}()

		data, err := io.ReadAll(resourceReader)
		if err != nil {
			return ast.WalkStop, errors.Wrapf(err, "could not read linked resource '%s'", destination)
		}

		mimeType := http.DetectContentType(data)

		if strings.HasSuffix(resourceURL.Path, ".svg") {
			mimeType = "image/svg+xml"
		}

		dataURL := dataurl.New(data, mimeType)

		image.Destination = []byte(dataURL.String())

		return ast.WalkContinue, nil
	})
	if err != nil {
		panic(errors.WithStack(err))
	}
}

var _ parser.ASTTransformer = &Transformer{}
