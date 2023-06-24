package tool

import (
	"context"

	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
)

type Options struct {
	Callbacks   []schema.Callback
	ParentRunID string
}

func Run(ctx context.Context, t schema.Tool, query string, optFns ...func(o *Options)) (string, error) {
	opts := Options{}

	for _, fn := range optFns {
		fn(&opts)
	}

	cm := callback.NewManager(opts.Callbacks, nil, false)

	rm, err := cm.OnToolStart(t.Name(), query)
	if err != nil {
		return "", err
	}

	output, err := t.Run(ctx, query)
	if err != nil {
		if cbErr := rm.OnToolError(err); cbErr != nil {
			return "", cbErr
		}

		return "", err
	}

	if err := rm.OnToolEnd(output); err != nil {
		return "", err
	}

	return output, nil
}
