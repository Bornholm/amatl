package dataurl

import (
	"context"
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

		ctx, err := pipeline.FromParserContext(pc)
		if err != nil {
			return ast.WalkStop, errors.WithStack(err)
		}

		switch typ := n.(type) {
		case *ast.Image:
			destination := string(typ.Destination)

			dataURL, err := t.toDataURL(ctx, destination)
			if err != nil {
				return ast.WalkStop, errors.WithStack(err)
			}

			typ.Destination = []byte(dataURL.String())

		default:
			return ast.WalkContinue, nil
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		panic(errors.WithStack(err))
	}
}

func (t *Transformer) toDataURL(ctx context.Context, destination string) (*dataurl.DataURL, error) {
	resourceURL, err := url.Parse(destination)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resourceReader, err := resolver.Resolve(ctx, resourceURL)
	if err != nil {
		return nil, errors.Wrapf(err, "could not resolve resource '%s'", destination)
	}

	defer func() {
		if err := resourceReader.Close(); err != nil {
			panic(errors.Wrapf(err, "could not close resource '%s'", resourceURL))
		}
	}()

	data, err := io.ReadAll(resourceReader)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read linked resource '%s'", destination)
	}

	mimeType := http.DetectContentType(data)

	if strings.HasSuffix(resourceURL.Path, ".svg") {
		mimeType = "image/svg+xml"
	}

	dataURL := dataurl.New(data, mimeType)

	return dataURL, nil
}

var _ parser.ASTTransformer = &Transformer{}

func isURL(rawURL string) bool {
	if _, err := url.ParseRequestURI(rawURL); err != nil {
		return false
	}

	return true
}
