package linkrewriter

import (
	"log/slog"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/Bornholm/amatl/pkg/pipeline"
	"github.com/Bornholm/amatl/pkg/resolver"
	"github.com/Bornholm/amatl/pkg/urlx"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type Transformer struct {
	replacements map[string]string
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

		workdir := resolver.ContextWorkDir(ctx)

		switch typ := n.(type) {
		case *ast.Link:
			destination := string(typ.Destination)

			if isURL(destination) || strings.HasPrefix(destination, "#") {
				return ast.WalkContinue, nil
			}

			if !filepath.IsAbs(destination) {
				destinationURL, err := urlx.Join(workdir, destination)
				if err != nil {
					return ast.WalkStop, errors.WithStack(err)
				}

				destination = destinationURL.String()
			}

			for prefix, replacement := range t.replacements {
				if strings.HasPrefix(destination, prefix) {
					updated := replacement + strings.TrimPrefix(destination, prefix)
					slog.DebugContext(ctx, "rewriting link", slog.String("original", destination), slog.String("updated", updated))
					destination = updated
					break
				}
			}

			typ.Destination = []byte(destination)

		default:
			return ast.WalkContinue, nil
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		panic(errors.WithStack(err))
	}
}

func NewTransformer(replacements map[string]string) *Transformer {
	return &Transformer{
		replacements: replacements,
	}
}

var _ parser.ASTTransformer = &Transformer{}

func isURL(str string) bool {
	_, err := url.ParseRequestURI(str)
	return err == nil
}
