package llm

import (
	"context"
	"fmt"
	"testing"

	huggingface "github.com/hupe1980/go-huggingface"
	"github.com/stretchr/testify/assert"
)

func TestHuggingFaceHub(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name                  string
		task                  string
		prompt                string
		expectedGeneratedText string
	}{
		{
			name:                  "Text Generation",
			task:                  "text-generation",
			prompt:                "Generate text",
			expectedGeneratedText: "Generate text",
		},
		{
			name:                  "Text2Text Generation",
			task:                  "text2text-generation",
			prompt:                "Convert text",
			expectedGeneratedText: "Converted text",
		},
		{
			name:                  "Summarization",
			task:                  "summarization",
			prompt:                "Summarize text",
			expectedGeneratedText: "Summarized text",
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock implementation of the HuggingFaceHubClient interface
			mockClient := &mockHuggingFaceHubClient{}

			// Create a new HuggingFaceHub instance with the updated options
			hub, err := NewHuggingFaceHubFromClient(mockClient, func(o *HuggingFaceHubOptions) {
				o.Task = tc.task
			})
			assert.NoError(t, err)

			// Mock the corresponding method of the client based on the task
			switch tc.task {
			case "text-generation":
				mockClient.textGenerationFn = func(ctx context.Context, req *huggingface.TextGenerationRequest) (huggingface.TextGenerationResponse, error) {
					return huggingface.TextGenerationResponse{{
						GeneratedText: fmt.Sprintf("%s%s", tc.expectedGeneratedText, tc.expectedGeneratedText),
					}}, nil
				}
			case "text2text-generation":
				mockClient.text2TextGenerationFn = func(ctx context.Context, req *huggingface.Text2TextGenerationRequest) (huggingface.Text2TextGenerationResponse, error) {
					return huggingface.Text2TextGenerationResponse{{
						GeneratedText: tc.expectedGeneratedText,
					}}, nil
				}
			case "summarization":
				mockClient.summarizationFn = func(ctx context.Context, req *huggingface.SummarizationRequest) (huggingface.SummarizationResponse, error) {
					return huggingface.SummarizationResponse{{
						SummaryText: tc.expectedGeneratedText,
					}}, nil
				}
			}

			// Generate text using the HuggingFaceHub instance
			result, err := hub.Generate(context.Background(), tc.prompt)
			assert.NoError(t, err)

			// Verify the generated text based on the task
			assert.Equal(t, tc.expectedGeneratedText, result.Generations[0].Text)
		})
	}

	t.Run("Type", func(t *testing.T) {
		// Create a hugging face hub instance
		llm, err := NewHuggingFaceHubFromClient(nil)
		assert.NoError(t, err)

		// Call the Type method
		typ := llm.Type()

		// Assert the result
		assert.Equal(t, "llm.HuggingFaceHub", typ)
	})

	t.Run("Verbose", func(t *testing.T) {
		// Create a hugging face hub instance
		llm, err := NewHuggingFaceHubFromClient(nil)
		assert.NoError(t, err)

		// Call the Verbose method
		verbose := llm.Verbose()

		// Assert the result
		assert.False(t, verbose)
	})

	t.Run("Callbacks", func(t *testing.T) {
		// Create a hugging face hub instance
		llm, err := NewHuggingFaceHubFromClient(nil)
		assert.NoError(t, err)

		// Call the Callbacks method
		callbacks := llm.Callbacks()

		// Assert the result
		assert.Empty(t, callbacks)
	})

	t.Run("InvocationParams", func(t *testing.T) {
		// Create a hugging face hub instance
		llm, err := NewHuggingFaceHubFromClient(&mockHuggingFaceHubClient{}, func(o *HuggingFaceHubOptions) {
			o.Model = "foo"
			o.Task = "bar"
		})
		assert.NoError(t, err)

		// Call the InvocationParams method
		params := llm.InvocationParams()

		// Assert the result
		assert.Equal(t, "foo", params["model"])
		assert.Equal(t, "bar", params["task"])
	})
}

type mockHuggingFaceHubClient struct {
	textGenerationFn      func(ctx context.Context, req *huggingface.TextGenerationRequest) (huggingface.TextGenerationResponse, error)
	text2TextGenerationFn func(ctx context.Context, req *huggingface.Text2TextGenerationRequest) (huggingface.Text2TextGenerationResponse, error)
	summarizationFn       func(ctx context.Context, req *huggingface.SummarizationRequest) (huggingface.SummarizationResponse, error)
	setModelFn            func(model string)
}

func (m *mockHuggingFaceHubClient) TextGeneration(ctx context.Context, req *huggingface.TextGenerationRequest) (huggingface.TextGenerationResponse, error) {
	if m.textGenerationFn != nil {
		return m.textGenerationFn(ctx, req)
	}

	return huggingface.TextGenerationResponse{}, nil
}

func (m *mockHuggingFaceHubClient) Text2TextGeneration(ctx context.Context, req *huggingface.Text2TextGenerationRequest) (huggingface.Text2TextGenerationResponse, error) {
	if m.text2TextGenerationFn != nil {
		return m.text2TextGenerationFn(ctx, req)
	}

	return huggingface.Text2TextGenerationResponse{}, nil
}

func (m *mockHuggingFaceHubClient) Summarization(ctx context.Context, req *huggingface.SummarizationRequest) (huggingface.SummarizationResponse, error) {
	if m.summarizationFn != nil {
		return m.summarizationFn(ctx, req)
	}

	return huggingface.SummarizationResponse{}, nil
}

func (m *mockHuggingFaceHubClient) SetModel(model string) {
	if m.setModelFn != nil {
		m.setModelFn(model)
	}
}
