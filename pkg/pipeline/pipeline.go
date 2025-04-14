package pipeline

import (
	"bytes"
	"context"
	"slices"
)

type Payload struct {
	attributes map[string]any
	data       []byte
}

func (p *Payload) SetAttribute(name string, value any) {
	p.attributes[name] = value
}

func (p *Payload) GetAttribute(name string) (any, bool) {
	value, exists := p.attributes[name]
	return value, exists
}

func (p *Payload) SetData(data []byte) {
	p.data = data
}

func (p *Payload) GetData() []byte {
	return p.data
}

func (p *Payload) Buffer() *bytes.Buffer {
	return bytes.NewBuffer(p.data)
}

func NewPayload(data []byte) *Payload {
	if data == nil {
		data = make([]byte, 0)
	}

	return &Payload{
		attributes: make(map[string]any),
		data:       data,
	}
}

type TransformerFunc func(ctx context.Context, payload *Payload) error

func (t TransformerFunc) Transform(ctx context.Context, payload *Payload) error {
	return t(ctx, payload)
}

var _ Transformer = TransformerFunc(func(ctx context.Context, payload *Payload) error { return nil })

type Transformer interface {
	Transform(ctx context.Context, payload *Payload) error
}

type Middleware func(next Transformer) Transformer

func Pipeline(middlewares ...Middleware) Transformer {
	slices.Reverse(middlewares)

	var transformer Transformer = TransformerFunc(func(ctx context.Context, payload *Payload) error {
		return nil
	})

	for _, m := range middlewares {
		transformer = m(transformer)
	}

	return transformer
}

func GetAttribute[T any](p *Payload, name string) (T, bool) {
	raw, ok := p.GetAttribute(name)
	if !ok {
		return *new(T), false
	}

	value, ok := raw.(T)
	if !ok {
		return *new(T), false
	}

	return value, true
}
