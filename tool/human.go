package tool

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Human satisfies the Tool interface.
var _ schema.Tool = (*Human)(nil)

// PromptFunc is a function type for displaying a prompt.
type PromptFunc = func(query string)

// InputFunc is a function type for retrieving user input.
type InputFunc = func() (string, error)

// HumanOptions contains options for configuring the Human tool.
type HumanOptions struct {
	// Function for displaying prompts.
	PromptFunc PromptFunc

	// Function for retrieving user input.
	InputFunc InputFunc
}

// Human is a tool that allows interaction with a human user.
type Human struct {
	opts HumanOptions
}

// NewHuman creates a new instance of the Human tool with the provided options.
func NewHuman(optFns ...func(o *HumanOptions)) *Human {
	opts := HumanOptions{
		PromptFunc: func(query string) {
			fmt.Println(query)
		},
		InputFunc: func() (string, error) {
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return "", err
			}
			input = strings.TrimSuffix(input, "\n")
			return input, nil
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Human{
		opts: opts,
	}
}

// Name returns the name of the tool.
func (t *Human) Name() string {
	return "Human"
}

// Description returns the description of the tool.
func (t *Human) Description() string {
	return `You can ask a human for guidance when you think you got stuck or you are not sure what to do next. The input should be a question for the human.`
}

// ArgsType returns the type of the input argument expected by the tool.
func (t *Human) ArgsType() reflect.Type {
	return reflect.TypeOf("") // string
}

// Run executes the tool with the given input and returns the output.
func (t *Human) Run(ctx context.Context, input any) (string, error) {
	query, ok := input.(string)
	if !ok {
		return "", errors.New("illegal input type")
	}

	t.opts.PromptFunc(query)

	return t.opts.InputFunc()
}

// Verbose returns the verbosity setting of the tool.
func (t *Human) Verbose() bool {
	return false
}

// Callbacks returns the registered callbacks of the tool.
func (t *Human) Callbacks() []schema.Callback {
	return nil
}
