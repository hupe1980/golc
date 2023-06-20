package tool

import (
	"context"
	"fmt"
	"time"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Sleep satisfies the Tool interface.
var _ schema.Tool = (*Sleep)(nil)

type Sleep struct {
	seconds int
}

func NewSleep(seconds int) *Sleep {
	return &Sleep{
		seconds: seconds,
	}
}

func (t *Sleep) Name() string {
	return "Sleep"
}

func (t *Sleep) Description() string {
	return `Make agent sleep for a specified number of seconds.`
}

func (t *Sleep) Run(ctx context.Context, query string) (string, error) {
	time.Sleep(time.Duration(t.seconds) * time.Second)
	return fmt.Sprintf("Agent slept for %d seconds.", t.seconds), nil
}
