package callback

import "github.com/hupe1980/golc/schema"

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

func (h *handler) OnLLMNewToken(token string) error {
	return nil
}

func (h *handler) OnLLMEnd(result *schema.LLMResult) error {
	return nil
}

func (h *handler) OnLLMError(llmError error) error {
	return nil
}

func (h *handler) OnChainStart(chainName string, inputs *schema.ChainValues) error {
	return nil
}

func (h *handler) OnChainEnd(outputs *schema.ChainValues) error {
	return nil
}

func (h *handler) OnChainError(chainError error) error {
	return nil
}
