package callback

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure StdOutHandler  satisfies the Callback interface.
var _ schema.Callback = (*StdOutHandler)(nil)

type StdOutHandler struct {
	handler
	writer io.Writer
}

func NewStdOutHandler() *StdOutHandler {
	return &StdOutHandler{
		writer: os.Stdout,
	}
}

func (cb *StdOutHandler) OnChainStart(ctx context.Context, input *schema.ChainStartInput) error {
	fmt.Fprintf(cb.writer, "\n\n\033[1m> Entering new %s chain...\033[0m\n", input.ChainType)
	return nil
}

func (cb *StdOutHandler) OnChainEnd(ctx context.Context, input *schema.ChainEndInput) error {
	fmt.Fprintln(cb.writer, "\n\033[1m> Finished chain.\033[0m")
	return nil
}

func (cb *StdOutHandler) OnAgentAction(ctx context.Context, input *schema.AgentActionInput) error {
	fmt.Fprintln(cb.writer, input.Action.Log)
	return nil
}

func (cb *StdOutHandler) OnAgentFinish(ctx context.Context, input *schema.AgentFinishInput) error {
	fmt.Fprintln(cb.writer, input.Finish.Log)
	return nil
}

func (cb *StdOutHandler) OnToolEnd(ctx context.Context, input *schema.ToolEndInput) error {
	fmt.Fprintln(cb.writer, input.Output)
	return nil
}

func (cb *StdOutHandler) OnText(ctx context.Context, input *schema.TextInput) error {
	fmt.Fprintln(cb.writer, input.Text)
	return nil
}
