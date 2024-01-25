package markdown

import "github.com/yuin/goldmark/ast"

// Render represents a single markdown rendering operation.
type Render struct {
	mr *Renderer

	emphToken   []byte
	strongToken []byte

	// TODO(bwplotka): Wrap it with something that catch errors.
	w      *LineIndentWriter
	source []byte
}

func (r *Render) Renderer() *Renderer {
	return r.mr
}

func (r *Render) Writer() *LineIndentWriter {
	return r.w
}

func (r *Render) EmphToken() []byte {
	return r.emphToken
}

func (r *Render) StrongToken() []byte {
	return r.strongToken
}

func (r *Render) Source() []byte {
	return r.source
}

func (r *Render) WrapNonEmptyContentWith(b []byte, entering bool) ast.WalkStatus {
	if entering {
		r.w.AddIndentOnFirstWrite(b)
		return ast.WalkContinue
	}

	if r.w.WasIndentOnFirstWriteWritten() {
		_, _ = r.w.Write(b)
		return ast.WalkContinue
	}
	r.w.DelIndentOnFirstWrite(b)
	return ast.WalkContinue
}
