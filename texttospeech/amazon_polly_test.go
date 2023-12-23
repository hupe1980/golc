package texttospeech

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/stretchr/testify/assert"
)

func TestAmazonPolly(t *testing.T) {
	tests := []struct {
		name          string
		client        AmazonPollyClient
		options       func(o *AmazonPollyOptions)
		inputText     string
		expectedError error
	}{
		{
			name: "SuccessfulSynthesis",
			client: &mockAmazonPollyClient{
				SynthesizeSpeechFn: func(ctx context.Context, params *polly.SynthesizeSpeechInput, optFns ...func(*polly.Options)) (*polly.SynthesizeSpeechOutput, error) {
					return &polly.SynthesizeSpeechOutput{
						AudioStream: &mockReadCloser{},
					}, nil
				},
			},
			inputText: "Hello, world!",
		},
		{
			name: "ClientError",
			client: &mockAmazonPollyClient{
				SynthesizeSpeechFn: func(ctx context.Context, params *polly.SynthesizeSpeechInput, optFns ...func(*polly.Options)) (*polly.SynthesizeSpeechOutput, error) {
					return nil, errors.New("mock client error")
				},
			},
			inputText:     "Error case",
			expectedError: errors.New("mock client error"),
		},
		{
			name: "UnsupportedResponseFormat",
			client: &mockAmazonPollyClient{
				SynthesizeSpeechFn: func(ctx context.Context, params *polly.SynthesizeSpeechInput, optFns ...func(*polly.Options)) (*polly.SynthesizeSpeechOutput, error) {
					return &polly.SynthesizeSpeechOutput{
						AudioStream: &mockReadCloser{},
					}, nil
				},
			},
			inputText: "Hello, world!",
			options: func(o *AmazonPollyOptions) {
				o.OutputFormat = "unsupported"
			},
			expectedError: errors.New("unsupported output format: unsupported"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := []func(o *AmazonPollyOptions){}
			if tt.options != nil {
				options = append(options, tt.options)
			}

			tts := NewAmazonPolly(tt.client, options...)

			audioStream, err := tts.SynthesizeSpeech(context.Background(), tt.inputText)

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

// mockAmazonPollyClient is a mock implementation of the AmazonPollyClient interface.
type mockAmazonPollyClient struct {
	SynthesizeSpeechFn func(ctx context.Context, params *polly.SynthesizeSpeechInput, optFns ...func(*polly.Options)) (*polly.SynthesizeSpeechOutput, error)
}

func (m *mockAmazonPollyClient) SynthesizeSpeech(ctx context.Context, params *polly.SynthesizeSpeechInput, optFns ...func(*polly.Options)) (*polly.SynthesizeSpeechOutput, error) {
	return m.SynthesizeSpeechFn(ctx, params, optFns...)
}
