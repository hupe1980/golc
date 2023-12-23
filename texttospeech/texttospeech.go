package texttospeech

import (
	"fmt"
	"io"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/vorbis"
	"github.com/hupe1980/golc/schema"
)

// audioStream is an implementation of the schema.AudioStream interface.
type audioStream struct {
	rc     io.ReadCloser
	format schema.OutputFormat
}

// NewAudioStream creates a new instance of the audioStream.
func NewAudioStream(rc io.ReadCloser, format schema.OutputFormat) schema.AudioStream {
	return &audioStream{
		rc:     rc,
		format: format,
	}
}

// Play plays the audio stream using the speaker package.
func (as *audioStream) Play() error {
	streamer, format, err := as.decode()
	if err != nil {
		return err
	}

	if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10)); err != nil {
		return err
	}

	done := make(chan struct{})

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		close(done)
	})))

	<-done

	return nil
}

// decode decodes the audio stream based on its format.
func (as *audioStream) decode() (s beep.StreamSeekCloser, format beep.Format, err error) {
	switch as.format {
	case schema.OutputFormatMP3:
		return mp3.Decode(as.rc)
	case schema.OutputFormatOggVorbis:
		return vorbis.Decode(as.rc)
	default:
		return nil, beep.Format{}, fmt.Errorf("unsupported output format: %s", as.format)
	}
}

// Save saves the audio stream to the provided writer.
func (as *audioStream) Save(dst io.Writer) error {
	_, err := io.Copy(dst, as.rc)
	if err != nil {
		return err
	}

	return nil
}

// Format returns the output format of the audio stream.
func (as *audioStream) Format() schema.OutputFormat {
	return as.format
}

// Read reads from the audio stream into the given byte slice.
func (as *audioStream) Read(p []byte) (n int, err error) {
	return as.rc.Read(p)
}

// Close closes the audio stream.
func (as *audioStream) Close() error {
	return as.rc.Close()
}
