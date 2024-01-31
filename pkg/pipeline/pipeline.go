package pipeline

import (
	"context"

	"github.com/pkg/errors"
)

type Pipeline struct {
	transformers []Transformer
}

// Transform implements Transformer.
func (p *Pipeline) Transform(ctx context.Context, input []byte) ([]byte, error) {
	var err error

	for _, tr := range p.transformers {
		input, err = tr.Transform(ctx, input)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return input, nil
}

func New(transformers ...Transformer) *Pipeline {
	return &Pipeline{
		transformers: transformers,
	}
}

var _ Transformer = &Pipeline{}
