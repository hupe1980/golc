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

type Sleep struct{}

func NewSleep() *Sleep {
	return &Sleep{}
}

func (t *Sleep) Name() string {
	return "Sleep"
}

func (t *Sleep) Description() string {
	return `Make agent sleep for a specified number of seconds.`
}

func (t *Sleep) ArgsType() reflect.Type {
	return reflect.TypeOf("") // string
}

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

func (t *Sleep) Verbose() bool {
	return false
}

func (t *Sleep) Callbacks() []schema.Callback {
	return nil
}
