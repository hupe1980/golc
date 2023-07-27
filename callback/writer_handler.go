package callback

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure WriterHandler satisfies the Callback interface.
var _ schema.Callback = (*WriterHandler)(nil)

type WriterHandlerOptions struct {
	Writer io.Writer
}

type WriterHandler struct {
	NoopHandler
	writer io.Writer
	opts   WriterHandlerOptions
}

func NewWriterHandler(optFns ...func(o *WriterHandlerOptions)) *WriterHandler {
	opts := WriterHandlerOptions{
		Writer: os.Stdout,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &WriterHandler{
		writer: opts.Writer,
		opts:   opts,
	}
}

func (cb *WriterHandler) OnChainStart(ctx context.Context, input *schema.ChainStartInput) error {
	fmt.Fprintf(cb.writer, "\n\n\033[1m> Entering new %s chain...\033[0m\n", input.ChainType)
	return nil
}

func (cb *WriterHandler) OnChainEnd(ctx context.Context, input *schema.ChainEndInput) error {
	fmt.Fprintln(cb.writer, "\n\033[1m> Finished chain.\033[0m")
	return nil
}

func (cb *WriterHandler) OnAgentAction(ctx context.Context, input *schema.AgentActionInput) error {
	fmt.Fprintln(cb.writer, input.Action.Log)
	return nil
}

func (cb *WriterHandler) OnAgentFinish(ctx context.Context, input *schema.AgentFinishInput) error {
	fmt.Fprintln(cb.writer, input.Finish.Log)
	return nil
}

func (cb *WriterHandler) OnToolEnd(ctx context.Context, input *schema.ToolEndInput) error {
	fmt.Fprintln(cb.writer, input.Output)
	return nil
}

func (cb *WriterHandler) OnText(ctx context.Context, input *schema.TextInput) error {
	fmt.Fprintln(cb.writer, input.Text)
	return nil
}
