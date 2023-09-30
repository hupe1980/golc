package llm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnforceStopTokens(t *testing.T) {
	tests := []struct {
		name string
		text string
		stop []string
		want string
	}{
		{
			name: "NoStopWords",
			text: "This is a test sentence.",
			stop: []string{"stop1", "stop2"},
			want: "This is a test sentence.",
		},
		{
			name: "StopWordsPresent",
			text: "Stop the text here.",
			stop: []string{"stop", "text"},
			want: "Stop the ",
		},
		{
			name: "EmptyText",
			text: "",
			stop: []string{"stop"},
			want: "",
		},
		{
			name: "EmptyStopWords",
			text: "This is a test sentence.",
			stop: []string{},
			want: "This is a test sentence.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EnforceStopTokens(tt.text, tt.stop)

			assert.Equal(t, tt.want, got, "unexpected result")
		})
	}
}
