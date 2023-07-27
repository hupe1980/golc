package callback

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure StreamWriterHandler satisfies the Callback interface.
var _ schema.Callback = (*StreamWriterHandler)(nil)

type StreamWriterHandlerOptions struct {
	Writer io.Writer
}

type StreamWriterHandler struct {
	NoopHandler
	writer io.Writer
	opts   StreamWriterHandlerOptions
}

func NewStreamWriterHandler(optFns ...func(o *StreamWriterHandlerOptions)) *StreamWriterHandler {
	opts := StreamWriterHandlerOptions{
		Writer: os.Stdout,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &StreamWriterHandler{
		writer: opts.Writer,
		opts:   opts,
	}
}

func (cb *StreamWriterHandler) AlwaysVerbose() bool {
	return true
}

func (cb *StreamWriterHandler) OnModelNewToken(ctx context.Context, input *schema.ModelNewTokenInput) error {
	fmt.Fprint(cb.writer, input.Token)
	return nil
}
