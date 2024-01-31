package pipeline

import "context"

type Transformer interface {
	Transform(ctx context.Context, input []byte) ([]byte, error)
}

type TransformFunc func(ctx context.Context, input []byte) ([]byte, error)

type transformer struct {
	transform TransformFunc
}

func (t *transformer) Transform(ctx context.Context, input []byte) ([]byte, error) {
	return t.transform(ctx, input)
}

func NewTransformer(transform TransformFunc) Transformer {
	return &transformer{
		transform: transform,
	}
}
