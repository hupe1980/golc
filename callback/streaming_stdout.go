package callback

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure StreamingStdOutHandler  satisfies the Callback interface.
var _ schema.Callback = (*StreamingStdOutHandler)(nil)

type StreamingStdOutHandler struct {
	handler
	writer io.Writer
}

func NewStreamingStdOutHandler() *StreamingStdOutHandler {
	return &StreamingStdOutHandler{
		writer: os.Stdout,
	}
}

func (cb *StreamingStdOutHandler) AlwaysVerbose() bool {
	return true
}

func (cb *StreamingStdOutHandler) OnModelNewToken(ctx context.Context, input *schema.ModelNewTokenInput) error {
	fmt.Fprint(cb.writer, input.Token)
	return nil
}
