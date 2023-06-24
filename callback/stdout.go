package callback

import (
	"fmt"
	"io"
	"os"

	"github.com/hupe1980/golc/schema"
)

type StdOutHandler struct {
	handler
	writer        io.Writer
	promptPrinter func(w io.Writer, llmName string, prompts []string) error
}

func NewStdOutHandler() *StdOutHandler {
	return &StdOutHandler{
		writer: os.Stdout,
		promptPrinter: func(w io.Writer, llmName string, prompts []string) error {
			for _, prompt := range prompts {
				fmt.Fprintln(w, prompt)
			}

			return nil
		},
	}
}

func (h *StdOutHandler) OnLLMStart(llmName string, prompts []string) error {
	return h.promptPrinter(h.writer, llmName, prompts)
}

func (h *StdOutHandler) OnChainStart(chainName string, inputs *schema.ChainValues) error {
	fmt.Fprintf(h.writer, "\n\n\033[1m> Entering new %s chain...\033[0m\n", chainName)

	return nil
}

func (h *StdOutHandler) OnChainEnd(outputs *schema.ChainValues) error {
	fmt.Fprintln(h.writer, "\n\033[1m> Finished chain.\033[0m")
	return nil
}

func (h *StdOutHandler) OnToolEnd(output string) error {
	//fmt.Println("\n\033[1m>XXXXX", output)
	return nil
}
