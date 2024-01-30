package layout

import "errors"

var (
	ErrSchemeNotRegistered = errors.New("scheme not registered")
	ErrTemplateNotFound    = errors.New("template not found")
)
