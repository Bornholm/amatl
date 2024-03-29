package node

import (
	"github.com/Bornholm/amatl/pkg/markdown/renderer/markdown"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
	extAST "github.com/yuin/goldmark/extension/ast"
)

type TaskCheckboxRenderer struct {
}

// Render implements NodeRenderer.
func (*TaskCheckboxRenderer) Render(r *markdown.Render, node ast.Node, entering bool) (ast.WalkStatus, error) {
	taskCheckbox, ok := node.(*extAST.TaskCheckBox)
	if !ok {
		return ast.WalkStop, errors.Errorf("expected *extAST.TaskCheckBox, got '%T'", node)
	}

	if !entering {
		return ast.WalkContinue, nil
	}

	if taskCheckbox.IsChecked {
		_, _ = r.Writer().Write([]byte("[X] "))
		return ast.WalkContinue, nil
	}

	_, _ = r.Writer().Write([]byte("[ ] "))

	return ast.WalkContinue, nil
}

var _ markdown.NodeRenderer = &TaskCheckboxRenderer{}
