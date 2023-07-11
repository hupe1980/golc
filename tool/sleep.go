package tool

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Sleep satisfies the Tool interface.
var _ schema.Tool = (*Sleep)(nil)

// Sleep is a tool that makes the agent sleep for a specified number of seconds.
type Sleep struct{}

// NewSleep creates a new instance of the Sleep tool.
func NewSleep() *Sleep {
	return &Sleep{}
}

// Name returns the name of the tool.
func (t *Sleep) Name() string {
	return "Sleep"
}

// Description returns the description of the tool.
func (t *Sleep) Description() string {
	return `Make agent sleep for a specified number of seconds.`
}

// ArgsType returns the type of the input argument expected by the tool.
func (t *Sleep) ArgsType() reflect.Type {
	return reflect.TypeOf("") // string
}

// Run executes the tool with the given input and returns the output.
func (t *Sleep) Run(ctx context.Context, input any) (string, error) {
	secondsStr, ok := input.(string)
	if !ok {
		return "", errors.New("illegal input type")
	}

	seconds, err := strconv.Atoi(secondsStr)
	if err != nil {
		return "", err
	}

	time.Sleep(time.Duration(seconds) * time.Second)

	return fmt.Sprintf("Agent slept for %d seconds.", seconds), nil
}

// Verbose returns the verbosity setting of the tool.
func (t *Sleep) Verbose() bool {
	return false
}

// Callbacks returns the registered callbacks of the tool.
func (t *Sleep) Callbacks() []schema.Callback {
	return nil
}
