package tool

import (
	"context"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Click satisfies the Tool interface.
var _ schema.Tool = (*Click)(nil)

type Click struct{}

func NewClick() *Click {
	return &Click{}
}

func (t *Click) Name() string {
	return "Click"
}

func (t *Click) Description() string {
	return `Click on an element with the given CSS selector.`
}

func (t *Click) Run(ctx context.Context, query string) (string, error) {
	return "", nil
}
