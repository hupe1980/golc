package chatmodel

import (
	"context"
	"testing"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestFake(t *testing.T) {
	t.Run("Generate_Success", func(t *testing.T) {
		// Arrange
		fake := NewSimpleFake("response")

		// Act
		result, err := fake.Generate(context.Background(), schema.ChatMessages{})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "response", result.Generations[0].Text)
	})

	t.Run("Generate_CustomResultFunc", func(t *testing.T) {
		// Arrange
		expectedResult := &schema.ModelResult{
			Generations: []schema.Generation{{Text: "custom_response"}},
		}
		fake := NewFake(func(ctx context.Context, messages schema.ChatMessages) (*schema.ModelResult, error) {
			return expectedResult, nil
		})

		// Act
		result, err := fake.Generate(context.Background(), schema.ChatMessages{})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "custom_response", result.Generations[0].Text)
	})

	t.Run("Generate_ErrorFromResultFunc", func(t *testing.T) {
		// Arrange
		expectedError := assert.AnError
		fake := NewFake(func(ctx context.Context, messages schema.ChatMessages) (*schema.ModelResult, error) {
			return nil, expectedError
		})

		// Act
		result, err := fake.Generate(context.Background(), schema.ChatMessages{})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, expectedError.Error())
	})

	t.Run("Type", func(t *testing.T) {
		// Arrange
		fake := NewSimpleFake("response")

		// Act
		result := fake.Type()

		// Assert
		assert.Equal(t, "chatmodel.Fake", result)
	})

	t.Run("Verbose", func(t *testing.T) {
		// Arrange
		fake := NewSimpleFake("response")

		// Act
		result := fake.Verbose()

		// Assert
		assert.Equal(t, golc.Verbose, result)
	})

	t.Run("Callbacks", func(t *testing.T) {
		// Arrange
		fake := NewSimpleFake("response")

		// Act
		result := fake.Callbacks()

		// Assert
		assert.Empty(t, result)
	})

	t.Run("InvocationParams", func(t *testing.T) {
		// Arrange
		fake := NewSimpleFake("response")

		// Act
		result := fake.InvocationParams()

		// Assert
		assert.Empty(t, result)
	})
}
