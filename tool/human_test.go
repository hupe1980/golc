package tool

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHuman(t *testing.T) {
	t.Run("Name", func(t *testing.T) {
		h := NewHuman()
		assert.Equal(t, "Human", h.Name())
	})

	t.Run("Description", func(t *testing.T) {
		h := NewHuman()
		expected := "You can ask a human for guidance when you think you got stuck or you are not sure what to do next. The input should be a question for the human."
		assert.Equal(t, expected, h.Description())
	})

	t.Run("ArgsType", func(t *testing.T) {
		h := NewHuman()
		expected := reflect.TypeOf("")
		assert.Equal(t, expected, h.ArgsType())
	})

	t.Run("Run", func(t *testing.T) {
		h := NewHuman()

		t.Run("ValidInput", func(t *testing.T) {
			ctx := context.Background()
			query := "What is your name?"
			expected := "John Doe"

			h.opts.InputFunc = func() (string, error) {
				return expected, nil
			}

			output, err := h.Run(ctx, query)
			assert.NoError(t, err)
			assert.Equal(t, expected, output)
		})

		t.Run("InvalidInputType", func(t *testing.T) {
			ctx := context.Background()
			query := 123 // Invalid input type

			_, err := h.Run(ctx, query)
			assert.Error(t, err)
			assert.Equal(t, "illegal input type", err.Error())
		})

		t.Run("ErrorRetrievingInput", func(t *testing.T) {
			ctx := context.Background()
			query := "What is your name?"

			expectedErr := errors.New("error retrieving input")
			h.opts.InputFunc = func() (string, error) {
				return "", expectedErr
			}

			_, err := h.Run(ctx, query)
			assert.Error(t, err)
			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("Verbose", func(t *testing.T) {
		h := NewHuman()
		assert.False(t, h.Verbose())
	})

	t.Run("Callbacks", func(t *testing.T) {
		h := NewHuman()
		assert.Nil(t, h.Callbacks())
	})
}
