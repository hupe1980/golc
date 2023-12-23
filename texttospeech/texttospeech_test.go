package texttospeech

import "io"

// mockReadCloser is a mock implementation of io.ReadCloser.
type mockReadCloser struct{}

func (m *mockReadCloser) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (m *mockReadCloser) Close() error {
	return nil
}
