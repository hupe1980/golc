package integration

import (
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

// Test case for messageTypeToOpenAIRole function
func TestMessageTypeToOpenAIRole(t *testing.T) {
	assertRole, assertErr := messageTypeToOpenAIRole(schema.ChatMessageTypeAI)
	assert.Equal(t, "assistant", assertRole)
	assert.NoError(t, assertErr)

	unknownRole, unknownErr := messageTypeToOpenAIRole("unknown")
	assert.Equal(t, "", unknownRole)
	assert.EqualError(t, unknownErr, "unknown message type: unknown")
}
