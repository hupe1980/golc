package tool

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSleep(t *testing.T) {
	sleepTool := NewSleep()

	// Test case for valid input
	t.Run("ValidInput", func(t *testing.T) {
		ctx := context.Background()
		input := "1" // 1 seconds

		expectedOutput := "Agent slept for 1 seconds."

		output, err := sleepTool.Run(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, output)
	})

	// Test case for invalid input type
	t.Run("InvalidInputType", func(t *testing.T) {
		ctx := context.Background()
		input := 10 // Invalid input type, expected string

		expectedOutput := ""
		expectedError := errors.New("illegal input type")

		output, err := sleepTool.Run(ctx, input)
		assert.Equal(t, expectedOutput, output)
		assert.EqualError(t, err, expectedError.Error())
	})

	// Test case for invalid input value (non-numeric)
	t.Run("InvalidInputValue", func(t *testing.T) {
		ctx := context.Background()
		input := "abc" // Invalid input value, expected numeric string

		expectedOutput := ""
		expectedError := errors.New("strconv.Atoi: parsing \"abc\": invalid syntax")

		output, err := sleepTool.Run(ctx, input)
		assert.Equal(t, expectedOutput, output)
		assert.EqualError(t, err, expectedError.Error())
	})
}
