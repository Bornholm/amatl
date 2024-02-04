package dataurl

import (
	"mime"
	"net/url"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/vincent-petithory/dataurl"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type Transformer struct {
	Cwd string
}

// Transform implements parser.ASTTransformer.
func (t *Transformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		image, ok := n.(*ast.Image)
		if !ok {
			return ast.WalkContinue, nil
		}

		destination := string(image.Destination)

		if isURL(destination) {
			return ast.WalkContinue, nil
		}

		if !filepath.IsAbs(destination) {
			destination = filepath.Join(t.Cwd, destination)
		}

		data, err := os.ReadFile(destination)
		if err != nil {
			return ast.WalkContinue, errors.Wrapf(err, "could not read linked resource '%s'", destination)
		}

		ext := filepath.Ext(destination)

		mimeType := mime.TypeByExtension(ext)
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		dataURL := dataurl.New(data, mimeType)

		image.Destination = []byte(dataURL.String())

		return ast.WalkContinue, nil
	})
}

var _ parser.ASTTransformer = &Transformer{}

func isURL(str string) bool {
	_, err := url.ParseRequestURI(str)
	return err == nil
}
