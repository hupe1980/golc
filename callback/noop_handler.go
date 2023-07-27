package callback

import (
	"context"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure NoopHandler satisfies the Callback interface.
var _ schema.Callback = (*NoopHandler)(nil)

type NoopHandler struct{}

func (h *NoopHandler) AlwaysVerbose() bool {
	return false
}

func (h *NoopHandler) RaiseError() bool {
	return false
}

func (h *NoopHandler) OnLLMStart(ctx context.Context, input *schema.LLMStartInput) error {
	return nil
}

func (h *NoopHandler) OnChatModelStart(ctx context.Context, input *schema.ChatModelStartInput) error {
	return nil
}

func (h *NoopHandler) OnModelNewToken(ctx context.Context, input *schema.ModelNewTokenInput) error {
	return nil
}

func (h *NoopHandler) OnModelEnd(ctx context.Context, input *schema.ModelEndInput) error {
	return nil
}

func (h *NoopHandler) OnModelError(ctx context.Context, input *schema.ModelErrorInput) error {
	return nil
}

func (h *NoopHandler) OnChainStart(ctx context.Context, input *schema.ChainStartInput) error {
	return nil
}

func (h *NoopHandler) OnChainEnd(ctx context.Context, input *schema.ChainEndInput) error {
	return nil
}

func (h *NoopHandler) OnChainError(ctx context.Context, input *schema.ChainErrorInput) error {
	return nil
}

func (h *NoopHandler) OnAgentAction(ctx context.Context, input *schema.AgentActionInput) error {
	return nil
}

func (h *NoopHandler) OnAgentFinish(ctx context.Context, input *schema.AgentFinishInput) error {
	return nil
}

func (h *NoopHandler) OnToolStart(ctx context.Context, input *schema.ToolStartInput) error {
	return nil
}
func (h *NoopHandler) OnToolEnd(ctx context.Context, input *schema.ToolEndInput) error {
	return nil
}
func (h *NoopHandler) OnToolError(ctx context.Context, input *schema.ToolErrorInput) error {
	return nil
}

func (h *NoopHandler) OnText(ctx context.Context, input *schema.TextInput) error {
	return nil
}

func (h *NoopHandler) OnRetrieverStart(ctx context.Context, input *schema.RetrieverStartInput) error {
	return nil
}

func (h *NoopHandler) OnRetrieverEnd(ctx context.Context, input *schema.RetrieverEndInput) error {
	return nil
}

func (h *NoopHandler) OnRetrieverError(ctx context.Context, input *schema.RetrieverErrorInput) error {
	return nil
}
