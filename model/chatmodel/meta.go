package chatmodel

import (
	"fmt"
	"strings"

	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
)

// convertMessagesToMetaPrompt converts a slice of chat messages into a formatted string suitable for Meta LLama Prompt.
// It handles different types of messages, such as generic chat messages, system messages, AI messages, and human messages.
// The resulting string includes role names, content, and specific formatting based on the message type.
// The function returns the formatted string and an error if any unsupported message type is encountered.
func convertMessagesToMetaPrompt(messages schema.ChatMessages) (string, error) {
	prompts := make([]string, len(messages))

	for i, message := range messages {
		switch v := message.(type) {
		case *schema.GenericChatMessage:
			prompts[i] = fmt.Sprintf("\n\n%s: %s", util.Capitalize(v.Role()), v.Content())
		case *schema.SystemChatMessage:
			prompts[i] = fmt.Sprintf("<<SYS>> %s <</SYS>>", v.Content())
		case *schema.AIChatMessage:
			prompts[i] = v.Content()
		case *schema.HumanChatMessage:
			prompts[i] = fmt.Sprintf("[INST] %s [/INST]", message.Content())
		default:
			return "", fmt.Errorf("unsupported message type: %s", message.Type())
		}
	}

	return strings.Join(prompts, "\n"), nil
}
