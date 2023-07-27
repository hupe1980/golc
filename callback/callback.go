// Package callback provides utilities for implementing callbacks.
package callback

import (
	"context"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure handler satisfies the Callback interface.
var _ schema.Callback = (*handler)(nil)

type handler struct{}

func (h *handler) AlwaysVerbose() bool {
	return false
}

func (h *handler) RaiseError() bool {
	return false
}

func (h *handler) OnLLMStart(ctx context.Context, input *schema.LLMStartInput) error {
	return nil
}

func (h *handler) OnChatModelStart(ctx context.Context, input *schema.ChatModelStartInput) error {
	return nil
}

func (h *handler) OnModelNewToken(ctx context.Context, input *schema.ModelNewTokenInput) error {
	return nil
}

func (h *handler) OnModelEnd(ctx context.Context, input *schema.ModelEndInput) error {
	return nil
}

func (h *handler) OnModelError(ctx context.Context, input *schema.ModelErrorInput) error {
	return nil
}

func (h *handler) OnChainStart(ctx context.Context, input *schema.ChainStartInput) error {
	return nil
}

func (h *handler) OnChainEnd(ctx context.Context, input *schema.ChainEndInput) error {
	return nil
}

func (h *handler) OnChainError(ctx context.Context, input *schema.ChainErrorInput) error {
	return nil
}

func (h *handler) OnAgentAction(ctx context.Context, input *schema.AgentActionInput) error {
	return nil
}

func (h *handler) OnAgentFinish(ctx context.Context, input *schema.AgentFinishInput) error {
	return nil
}

func (h *handler) OnToolStart(ctx context.Context, input *schema.ToolStartInput) error {
	return nil
}
func (h *handler) OnToolEnd(ctx context.Context, input *schema.ToolEndInput) error {
	return nil
}
func (h *handler) OnToolError(ctx context.Context, input *schema.ToolErrorInput) error {
	return nil
}

func (h *handler) OnText(ctx context.Context, input *schema.TextInput) error {
	return nil
}

func (h *handler) OnRetrieverStart(ctx context.Context, input *schema.RetrieverStartInput) error {
	return nil
}

func (h *handler) OnRetrieverEnd(ctx context.Context, input *schema.RetrieverEndInput) error {
	return nil
}

func (h *handler) OnRetrieverError(ctx context.Context, input *schema.RetrieverErrorInput) error {
	return nil
}
