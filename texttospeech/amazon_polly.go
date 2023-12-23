package texttospeech

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/polly/types"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure AmazonPolly satisfies the TextToSpeech interface.
var _ schema.TextToSpeech = (*AmazonPolly)(nil)

// AmazonPollyClient is an interface for the Amazon Polly Text-to-Speech API client.
type AmazonPollyClient interface {
	SynthesizeSpeech(ctx context.Context, params *polly.SynthesizeSpeechInput, optFns ...func(*polly.Options)) (*polly.SynthesizeSpeechOutput, error)
}

// AmazonPollyOptions contains options for configuring the Amazon Polly transformer.
type AmazonPollyOptions struct {
	VoiceID         types.VoiceId
	Engine          types.Engine
	OutputFormat    types.OutputFormat
	LanguageCode    types.LanguageCode
	LexiconNames    []string
	SampleRate      string
	SpeechMarkTypes []types.SpeechMarkType
	TextType        types.TextType
}

// AmazonPolly is a transformer that uses the Amazon Polly Text-to-Speech API to synthesize speech from text.
type AmazonPolly struct {
	client AmazonPollyClient
	opts   AmazonPollyOptions
}

// NewAmazonPolly creates a new instance of the Amazon Polly transformer.
func NewAmazonPolly(client AmazonPollyClient, optFns ...func(o *AmazonPollyOptions)) *AmazonPolly {
	opts := AmazonPollyOptions{
		VoiceID:      types.VoiceIdJoanna,
		Engine:       types.EngineStandard,
		OutputFormat: types.OutputFormatMp3,
		SampleRate:   "22050",
		TextType:     types.TextTypeText,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &AmazonPolly{
		client: client,
		opts:   opts,
	}
}

// SynthesizeSpeech uses the Amazon Polly Text-to-Speech API to transform text into an audio stream.
func (t2s *AmazonPolly) SynthesizeSpeech(ctx context.Context, text string) (schema.AudioStream, error) {
	var internalOutputFormat schema.OutputFormat

	switch t2s.opts.OutputFormat {
	case types.OutputFormatMp3:
		internalOutputFormat = schema.OutputFormatMP3
	case types.OutputFormatOggVorbis:
		internalOutputFormat = schema.OutputFormatOggVorbis
	case types.OutputFormatPcm:
		internalOutputFormat = schema.OutputFormatPCM
	default:
		return nil, fmt.Errorf("unsupported output format: %s", t2s.opts.OutputFormat)
	}

	res, err := t2s.client.SynthesizeSpeech(ctx, &polly.SynthesizeSpeechInput{
		Text:            aws.String(text),
		OutputFormat:    t2s.opts.OutputFormat,
		VoiceId:         t2s.opts.VoiceID,
		Engine:          t2s.opts.Engine,
		LanguageCode:    t2s.opts.LanguageCode,
		LexiconNames:    t2s.opts.LexiconNames,
		SampleRate:      aws.String(t2s.opts.SampleRate),
		SpeechMarkTypes: t2s.opts.SpeechMarkTypes,
		TextType:        t2s.opts.TextType,
	})
	if err != nil {
		return nil, err
	}

	return NewAudioStream(res.AudioStream, internalOutputFormat), nil
}
