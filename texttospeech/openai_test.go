package texttospeech

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

// MockOpenAIClient is a mock implementation of the OpenAIClient interface.
type MockOpenAIClient struct {
	CreateSpeechFn func(ctx context.Context, request openai.CreateSpeechRequest) (response io.ReadCloser, err error)
}

func (m *MockOpenAIClient) CreateSpeech(ctx context.Context, request openai.CreateSpeechRequest) (response io.ReadCloser, err error) {
	return m.CreateSpeechFn(ctx, request)
}

func TestOpenAI(t *testing.T) {
	tests := []struct {
		name          string
		client        OpenAIClient
		options       func(o *OpenAIOptions)
		inputText     string
		expectedError error
	}{
		{
			name: "SuccessfulSynthesis",
			client: &MockOpenAIClient{
				CreateSpeechFn: func(ctx context.Context, request openai.CreateSpeechRequest) (response io.ReadCloser, err error) {
					return &mockReadCloser{}, nil
				},
			},
			inputText: "Hello, world!",
		},
		{
			name: "ClientError",
			client: &MockOpenAIClient{
				CreateSpeechFn: func(ctx context.Context, request openai.CreateSpeechRequest) (response io.ReadCloser, err error) {
					return nil, errors.New("mock client error")
				},
			},
			inputText:     "Error case",
			expectedError: errors.New("mock client error"),
		},
		{
			name: "UnsupportedResponseFormat",
			client: &MockOpenAIClient{
				CreateSpeechFn: func(ctx context.Context, request openai.CreateSpeechRequest) (response io.ReadCloser, err error) {
					return &mockReadCloser{}, nil
				},
			},
			options: func(o *OpenAIOptions) {
				o.ResponseFormat = "unsupported"
			},
			inputText:     "Unsupported format",
			expectedError: errors.New("unsupported response format: unsupported"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := []func(o *OpenAIOptions){}
			if tt.options != nil {
				options = append(options, tt.options)
			}

			t2s := NewOpenAIFromClient(tt.client, options...)

			audioStream, err := t2s.SynthesizeSpeech(context.Background(), tt.inputText)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error(), "Unexpected error")
				assert.Nil(t, audioStream, "AudioStream should be nil on error")
			} else {
				assert.NoError(t, err, "Unexpected error")
				assert.NotNil(t, audioStream, "AudioStream should not be nil on success")
			}
		})
	}
}
