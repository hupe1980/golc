package callback

import (
	"fmt"
	"io"
	"os"

	"github.com/hupe1980/golc/schema"
)

type StdOutHandler struct {
	handler
	writer io.Writer
}

func NewStdOutHandler() *StdOutHandler {
	return &StdOutHandler{
		writer: os.Stdout,
	}
}

func (h *StdOutHandler) OnChainStart(chainName string, inputs schema.ChainValues) error {
	fmt.Fprintf(h.writer, "\n\n\033[1m> Entering new %s chain...\033[0m\n", chainName)
	return nil
}

func (h *StdOutHandler) OnChainEnd(outputs schema.ChainValues) error {
	fmt.Fprintln(h.writer, "\n\033[1m> Finished chain.\033[0m")
	return nil
}

func (h *StdOutHandler) OnAgentAction(action schema.AgentAction) error {
	fmt.Fprintln(h.writer, action.Log)
	return nil
}

func (h *StdOutHandler) OnAgentFinish(finish schema.AgentFinish) error {
	fmt.Fprintln(h.writer, finish.Log)
	return nil
}

func (h *StdOutHandler) OnToolEnd(output string) error {
	fmt.Fprintln(h.writer, output)
	return nil
}

func (h *StdOutHandler) OnText(text string) error {
	fmt.Fprintln(h.writer, text)
	return nil
}
