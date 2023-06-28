package callback

import "github.com/hupe1980/golc/schema"

// Compile time check to ensure handler satisfies the Callback interface.
var _ schema.Callback = (*handler)(nil)

type handler struct{}

func (h *handler) AlwaysVerbose() bool {
	return false
}

func (h *handler) RaiseError() bool {
	return false
}

func (h *handler) OnLLMStart(llmName string, prompts []string) error {
	return nil
}

func (h *handler) OnChatModelStart(chatModelName string, messages []schema.ChatMessages) error {
	return nil
}

func (h *handler) OnModelNewToken(token string) error {
	return nil
}

func (h *handler) OnModelEnd(result schema.LLMResult) error {
	return nil
}

func (h *handler) OnModelError(llmError error) error {
	return nil
}

func (h *handler) OnChainStart(chainName string, inputs schema.ChainValues) error {
	return nil
}

func (h *handler) OnChainEnd(outputs schema.ChainValues) error {
	return nil
}

func (h *handler) OnChainError(chainError error) error {
	return nil
}

func (h *handler) OnAgentAction(action schema.AgentAction) error {
	return nil
}

func (h *handler) OnAgentFinish(finish schema.AgentFinish) error {
	return nil
}

func (h *handler) OnToolStart(toolName string, input string) error {
	return nil
}
func (h *handler) OnToolEnd(output string) error {
	return nil
}
func (h *handler) OnToolError(toolError error) error {
	return nil
}

func (h *handler) OnText(text string) error {
	return nil
}
