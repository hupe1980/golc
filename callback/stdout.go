package callback

import (
	"fmt"
	"io"
	"os"

	"github.com/hupe1980/golc"
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

func (h *StdOutHandler) OnChainStart(chainName string, inputs *golc.ChainValues) error {
	fmt.Fprintf(h.writer, "\n\n\033[1m> Entering new %s chain...\033[0m\n", chainName)

	return nil
}

func (h *StdOutHandler) OnChainEnd(outputs *golc.ChainValues) error {
	fmt.Fprintln(h.writer, "\n\033[1m> Finished chain.\033[0m")
	return nil
}
