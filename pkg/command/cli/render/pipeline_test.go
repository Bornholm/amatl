package render

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type transformerFunc func()

func (fn transformerFunc) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	fn()
}

func parseWithTransformer(fn func()) (ast.Node, error) {
	parse := goldmark.New().Parser()
	parse.AddOptions(parser.WithASTTransformers(util.Prioritized(transformerFunc(fn), 0)))

	return parseDocument(parse, text.NewReader([]byte("text")), parser.NewContext())
}

func TestParseDocumentReturnsError(t *testing.T) {
	_, err := parseWithTransformer(func() {
		panic(errors.New("boom"))
	})
	if err == nil || err.Error() != "boom" {
		t.Errorf("parseDocument: expected the 'boom' error, got %v", err)
	}
}

func TestParseDocumentRepanicsRuntimeError(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("parseDocument: expected a runtime error to panic")
		}
	}()

	var nodes []ast.Node

	_, _ = parseWithTransformer(func() {
		_ = nodes[1]
	})
}
