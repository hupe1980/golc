package texttospeech

import (
	"context"
	"fmt"

	"github.com/hupe1980/golc/schema"
	"github.com/sashabaranov/go-openai"
)

// Compile time check to ensure OpenAI satisfies the TextToSpeech interface.
var _ schema.TextToSpeech = (*OpenAI)(nil)

// OpenAIClient is an interface for the OpenAI Text-to-Speech API client.
type OpenAIClient interface {
	CreateSpeech(ctx context.Context, request openai.CreateSpeechRequest) (response openai.RawResponse, err error)
}

// OpenAIOptions contains options for configuring the OpenAI transformer.
type OpenAIOptions struct {
	Model          openai.SpeechModel
	Voice          openai.SpeechVoice
	ResponseFormat openai.SpeechResponseFormat
	Speed          float64
}

// DefaultOpenAIOptions provides default values for OpenAIOptions.
var DefaultOpenAIOptions = OpenAIOptions{
	Model:          openai.TTSModel1,
	Voice:          openai.VoiceAlloy,
	ResponseFormat: openai.SpeechResponseFormatMp3,
	Speed:          1.0,
}

// OpenAI is a transformer that uses the OpenAI Text-to-Speech API to synthesize speech from text.
type OpenAI struct {
	client OpenAIClient
	opts   OpenAIOptions
}

// NewOpenAI creates a new instance of the OpenAI transformer.
func NewOpenAI(apiKey string, optFns ...func(o *OpenAIOptions)) *OpenAI {
	opts := DefaultOpenAIOptions

	for _, fn := range optFns {
		fn(&opts)
	}

	config := openai.DefaultConfig(apiKey)

	client := openai.NewClientWithConfig(config)

	return NewOpenAIFromClient(client, func(o *OpenAIOptions) {
		*o = opts
	})
}

// NewOpenAI creates a new instance of the OpenAI transformer.
func NewOpenAIFromClient(client OpenAIClient, optFns ...func(o *OpenAIOptions)) *OpenAI {
	opts := DefaultOpenAIOptions

	for _, fn := range optFns {
		fn(&opts)
	}

	return &OpenAI{
		client: client,
		opts:   opts,
	}
}

// SynthesizeSpeech uses the OpenAI Text-to-Speech API to transform text into an audio stream.
func (t2s *OpenAI) SynthesizeSpeech(ctx context.Context, text string) (schema.AudioStream, error) {
	var internalOutputFormat schema.OutputFormat

	switch t2s.opts.ResponseFormat {
	case openai.SpeechResponseFormatMp3:
		internalOutputFormat = schema.OutputFormatMP3
	case openai.SpeechResponseFormatAac:
		internalOutputFormat = schema.OutputFormatAAC
	case openai.SpeechResponseFormatOpus:
		internalOutputFormat = schema.OutputFormatOpus
	case openai.SpeechResponseFormatFlac:
		internalOutputFormat = schema.OutputFormatFlac
	default:
		return nil, fmt.Errorf("unsupported response format: %s", t2s.opts.ResponseFormat)
	}

	res, err := t2s.client.CreateSpeech(ctx, openai.CreateSpeechRequest{
		Model:          t2s.opts.Model,
		Input:          text,
		Voice:          t2s.opts.Voice,
		ResponseFormat: t2s.opts.ResponseFormat,
		Speed:          t2s.opts.Speed,
	})
	if err != nil {
		return nil, err
	}

	return NewAudioStream(res, internalOutputFormat), nil
}
