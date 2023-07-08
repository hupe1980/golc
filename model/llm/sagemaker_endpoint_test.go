package llm

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntime"
	"github.com/stretchr/testify/assert"
)

func TestSagemakerEndpoint(t *testing.T) {
	t.Run("Generate", func(t *testing.T) {
		t.Run("Successful generation", func(t *testing.T) {
			mockTransformer := &mockTransformer{}
			mockSagemakerClient := &mockSagemakerClient{}

			contentHandler := NewLLMContentHandler("text/plain", "text/plain", mockTransformer)

			endpoint, err := NewSagemakerEndpoint(mockSagemakerClient, "my-endpoint", contentHandler)
			assert.NoError(t, err)

			expectedPrompt := "Hello, world!"
			expectedOutput := "Generated text"

			mockSagemakerClient.InvokeEndpointFunc = func(ctx context.Context, params *sagemakerruntime.InvokeEndpointInput, optFns ...func(*sagemakerruntime.Options)) (*sagemakerruntime.InvokeEndpointOutput, error) {
				assert.Equal(t, expectedPrompt, string(params.Body))
				assert.Equal(t, aws.String(contentHandler.ContentType()), params.ContentType)
				assert.Equal(t, aws.String(contentHandler.Accept()), params.Accept)
				return &sagemakerruntime.InvokeEndpointOutput{
					Body: []byte(expectedOutput),
				}, nil
			}

			result, err := endpoint.Generate(context.Background(), expectedPrompt)
			assert.NoError(t, err)
			assert.Len(t, result.Generations, 1)
			assert.Equal(t, expectedOutput, result.Generations[0].Text)
			assert.Empty(t, result.LLMOutput)
		})

		t.Run("Failed content transformation", func(t *testing.T) {
			mockTransformer := &mockTransformer{}
			mockSagemakerClient := &mockSagemakerClient{}

			contentHandler := NewLLMContentHandler("text/plain", "text/plain", mockTransformer)

			endpoint, err := NewSagemakerEndpoint(mockSagemakerClient, "my-endpoint", contentHandler)
			assert.NoError(t, err)

			expectedError := errors.New("Content transformation error")

			mockTransformer.TransformInputError = expectedError

			result, err := endpoint.Generate(context.Background(), "Invalid prompt")
			assert.Error(t, err)
			assert.Equal(t, err, expectedError)
			assert.Nil(t, result)
		})

		t.Run("Failed endpoint invocation", func(t *testing.T) {
			mockTransformer := &mockTransformer{}
			mockSagemakerClient := &mockSagemakerClient{}

			contentHandler := NewLLMContentHandler("text/plain", "text/plain", mockTransformer)

			endpoint, err := NewSagemakerEndpoint(mockSagemakerClient, "my-endpoint", contentHandler)
			assert.NoError(t, err)

			expectedError := errors.New("Endpoint invocation error")

			mockSagemakerClient.InvokeEndpointFunc = func(ctx context.Context, params *sagemakerruntime.InvokeEndpointInput, optFns ...func(*sagemakerruntime.Options)) (*sagemakerruntime.InvokeEndpointOutput, error) {
				return nil, expectedError
			}

			result, err := endpoint.Generate(context.Background(), "Hello, world!")
			assert.Error(t, err)
			assert.Equal(t, err, expectedError)
			assert.Nil(t, result)
		})
	})

	t.Run("Type", func(t *testing.T) {
		endpoint, err := NewSagemakerEndpoint(nil, "", nil)
		assert.NoError(t, err)
		assert.Equal(t, "llm.SagemakerEndpoint", endpoint.Type())
	})

	t.Run("Verbose", func(t *testing.T) {
		endpoint, err := NewSagemakerEndpoint(nil, "", nil)
		assert.NoError(t, err)
		assert.False(t, endpoint.Verbose())
	})

	t.Run("Callbacks", func(t *testing.T) {
		endpoint, err := NewSagemakerEndpoint(nil, "", nil)
		assert.NoError(t, err)
		// Call the Callbacks method
		callbacks := endpoint.Callbacks()

		// Assert the result
		assert.Empty(t, callbacks)
	})

	t.Run("InvocationParams", func(t *testing.T) {
		endpoint, err := NewSagemakerEndpoint(nil, "", nil)
		assert.NoError(t, err)
		params := endpoint.InvocationParams()
		assert.Empty(t, params)
	})
}

type mockTransformer struct {
	TransformInputError  error
	TransformOutputError error
}

func (mt *mockTransformer) TransformInput(prompt string) ([]byte, error) {
	if mt.TransformInputError != nil {
		return nil, mt.TransformInputError
	}

	return []byte(prompt), nil
}

func (mt *mockTransformer) TransformOutput(output []byte) (string, error) {
	if mt.TransformOutputError != nil {
		return "", mt.TransformOutputError
	}

	return string(output), nil
}

type mockSagemakerClient struct {
	InvokeEndpointFunc func(ctx context.Context, params *sagemakerruntime.InvokeEndpointInput, optFns ...func(*sagemakerruntime.Options)) (*sagemakerruntime.InvokeEndpointOutput, error)
}

func (m *mockSagemakerClient) InvokeEndpoint(ctx context.Context, params *sagemakerruntime.InvokeEndpointInput, optFns ...func(*sagemakerruntime.Options)) (*sagemakerruntime.InvokeEndpointOutput, error) {
	return m.InvokeEndpointFunc(ctx, params, optFns...)
}
