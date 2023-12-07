package chatmodel

import (
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestConvertMessagesToMetaPrompt(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		inputMessages  schema.ChatMessages
		expectedResult string
		expectedError  error
	}{
		{
			name: "GenericChatMessage",
			inputMessages: schema.ChatMessages{
				schema.NewGenericChatMessage("Hello!", "user"),
			},
			expectedResult: "\n\nUser: Hello!",
			expectedError:  nil,
		},
		{
			name: "SystemChatMessage",
			inputMessages: schema.ChatMessages{
				schema.NewSystemChatMessage("System message"),
			},
			expectedResult: "<<SYS>> System message <</SYS>>",
			expectedError:  nil,
		},
		{
			name: "AIChatMessage",
			inputMessages: schema.ChatMessages{
				schema.NewAIChatMessage("AI response"),
			},
			expectedResult: "AI response",
			expectedError:  nil,
		},
		{
			name: "HumanChatMessage",
			inputMessages: schema.ChatMessages{
				schema.NewHumanChatMessage("User message"),
			},
			expectedResult: "[INST] User message [/INST]",
			expectedError:  nil,
		},
		{
			name: "Conversation",
			inputMessages: schema.ChatMessages{
				schema.NewSystemChatMessage("You're an assistant"),
				schema.NewHumanChatMessage("Hello"),
				schema.NewAIChatMessage("Answer:"),
			},
			expectedResult: "<<SYS>> You're an assistant <</SYS>>\n[INST] Hello [/INST]\nAnswer:",
			expectedError:  nil,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function under test
			result, err := convertMessagesToMetaPrompt(tt.inputMessages)

			// Check the results
			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
