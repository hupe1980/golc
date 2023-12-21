package stream

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStream(t *testing.T) {
	t.Run("Recv - Successful JSON Unmarshal", func(t *testing.T) {
		responseBody := `{"field": "value"}`
		reader := &mockReader{Data: []byte(responseBody)}
		response := &http.Response{Body: io.NopCloser(reader)}

		stream := NewStream[mockResponse](response)
		defer stream.Close()

		res, err := stream.Recv()
		assert.NoError(t, err)
		assert.Equal(t, "value", res.Field)
	})

	t.Run("Recv - Error in JSON Unmarshal", func(t *testing.T) {
		responseBody := `invalid JSON`
		reader := &mockReader{Data: []byte(responseBody)}
		response := &http.Response{Body: io.NopCloser(reader)}

		stream := NewStream[mockResponse](response)
		defer stream.Close()

		res, err := stream.Recv()
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Recv - Error Reading from Stream", func(t *testing.T) {
		expectedErr := errors.New("test error")
		reader := &mockReader{Err: expectedErr}
		response := &http.Response{Body: io.NopCloser(reader)}

		stream := NewStream[mockResponse](response)
		defer stream.Close()

		res, err := stream.Recv()
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, res)
	})

	t.Run("Close - Close Stream Successfully", func(t *testing.T) {
		reader := &mockReader{Data: []byte{}}
		response := &http.Response{Body: io.NopCloser(reader)}

		stream := NewStream[mockResponse](response)
		err := stream.Close()
		assert.NoError(t, err)
	})
}

type mockResponse struct {
	Field string `json:"field"`
}

// mockReader is a mock implementation of io.Reader for testing purposes.
type mockReader struct {
	Data []byte
	Err  error
}

func (m *mockReader) Read(p []byte) (n int, err error) {
	if m.Err != nil {
		return 0, m.Err
	}

	l := copy(p, m.Data)

	return l, io.EOF
}
