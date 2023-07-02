package callback

import (
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

func (cb *StdOutHandler) OnChainStart(chainName string, inputs schema.ChainValues) error {
	fmt.Fprintf(cb.writer, "\n\n\033[1m> Entering new %s chain...\033[0m\n", chainName)
	return nil
}

func (cb *StdOutHandler) OnChainEnd(outputs schema.ChainValues) error {
	fmt.Fprintln(cb.writer, "\n\033[1m> Finished chain.\033[0m")
	return nil
}

func (cb *StdOutHandler) OnAgentAction(action schema.AgentAction) error {
	fmt.Fprintln(cb.writer, action.Log)
	return nil
}

func (cb *StdOutHandler) OnAgentFinish(finish schema.AgentFinish) error {
	fmt.Fprintln(cb.writer, finish.Log)
	return nil
}

func (cb *StdOutHandler) OnToolEnd(output string) error {
	fmt.Fprintln(cb.writer, output)
	return nil
}

func (cb *StdOutHandler) OnText(text string) error {
	fmt.Fprintln(cb.writer, text)
	return nil
}
