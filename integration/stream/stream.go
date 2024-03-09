package stream

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// StreamOptions contains options for configuring the Stream.
type StreamOptions struct {
	Delimeter []byte
}

// Stream is a generic streaming reader for decoding JSON-encoded data from an HTTP response.
type Stream[T any] struct {
	reader streamReader
	closer io.Closer
}

// NewStream creates a new instance of the Stream.
func NewStream[T any](response *http.Response, optFns ...func(o *StreamOptions)) *Stream[T] {
	opts := StreamOptions{}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Stream[T]{
		reader: newStreamReader(response.Body, opts.Delimeter),
		closer: response.Body,
	}
}

// Recv reads and decodes the next value from the stream.
func (s Stream[T]) Recv() (*T, error) {
	value := new(T)

	bytes, err := s.reader.ReadFromStream()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, value); err != nil {
		return nil, err
	}

	return value, nil
}

// Close closes the underlying stream.
func (s Stream[T]) Close() error {
	return s.closer.Close()
}

// streamReader is an interface for reading from a stream.
type streamReader interface {
	ReadFromStream() ([]byte, error)
}

// newStreamReader creates a new instance of streamReader based on the provided delimiter.
func newStreamReader(reader io.Reader, delimiter []byte) streamReader {
	return newScannerStreamReader(reader, delimiter)
}

// scannerStreamReader is a streamReader implementation using a bufio.Scanner with a custom delimiter.
type scannerStreamReader struct {
	scanner *bufio.Scanner
}

// newScannerStreamReader creates a new instance of scannerStreamReader with a custom delimiter.
func newScannerStreamReader(reader io.Reader, delimiter []byte) *scannerStreamReader {
	scanner := bufio.NewScanner(reader)

	if delimiter != nil {
		scanner.Split(func(data []byte, atEOF bool) (int, []byte, error) {
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}

			if i := bytes.Index(data, delimiter); i >= 0 {
				return i + len(delimiter), data[0:i], nil
			}

			if atEOF {
				return len(data), data, nil
			}

			return 0, nil, nil
		})
	}

	return &scannerStreamReader{
		scanner: scanner,
	}
}

// ReadFromStream reads a chunk from the stream using a bufio.Scanner with a custom delimiter.
func (b *scannerStreamReader) ReadFromStream() ([]byte, error) {
	if b.scanner.Scan() {
		return b.scanner.Bytes(), nil
	}

	if err := b.scanner.Err(); err != nil {
		return nil, err
	}

	return nil, io.EOF
}
