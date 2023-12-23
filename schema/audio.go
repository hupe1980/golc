package schema

import (
	"context"
	"io"
)

// OutputFormat defines the supported audio output formats.
type OutputFormat string

const (
	// OutputFormatMP3 represents the MP3 audio format.
	OutputFormatMP3 OutputFormat = "mp3"
	// OutputFormatOggVorbis represents the Ogg Vorbis audio format.
	OutputFormatOggVorbis OutputFormat = "ogg_vorbis"
	// OutputFormatPCM represents the PCM audio format.
	OutputFormatPCM OutputFormat = "pcm"
	// OutputFormatOpus represents the Opus audio format.
	OutputFormatOpus OutputFormat = "opus"
	// OutputFormatAAC represents the AAC audio format.
	OutputFormatAAC OutputFormat = "aac"
	// OutputFormatFlac represents the FLAC audio format.
	OutputFormatFlac OutputFormat = "flac"
)

// AudioStream is an interface for handling audio streams.
type AudioStream interface {
	// Play plays the audio stream.
	Play() error
	// Save saves the audio stream to the provided writer.
	Save(dst io.Writer) error
	// Format returns the output format of the audio stream.
	Format() OutputFormat
	// Read reads from the audio stream into the given byte slice.
	Read(p []byte) (n int, err error)
	// Close closes the audio stream.
	Close() error
}

// TextToSpeech is an interface for converting text to speech.
type TextToSpeech interface {
	// SynthesizeSpeech converts the given text to an audio stream.
	SynthesizeSpeech(ctx context.Context, text string) (AudioStream, error)
}
