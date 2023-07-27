package callback

import (
	"context"

	"github.com/google/uuid"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure manager satisfies the CallbackManager interface.
var _ schema.CallbackManager = (*manager)(nil)

// Compile time check to ensure manager satisfies the CallbackManagerForChainRun interface.
var _ schema.CallbackManagerForChainRun = (*manager)(nil)

// Compile time check to ensure manager satisfies the CallbackManagerForModelRun interface.
var _ schema.CallbackManagerForModelRun = (*manager)(nil)

// Compile time check to ensure manager satisfies the CallbackManagerForToolRun interface.
var _ schema.CallbackManagerForToolRun = (*manager)(nil)

// Compile time check to ensure manager satisfies the CallbackManagerForRetrieverRun interface.
var _ schema.CallbackManagerForRetrieverRun = (*manager)(nil)

type ManagerOptions struct {
	ParentRunID string
}

type manager struct {
	callbacks            []schema.Callback
	inheritableCallbacks []schema.Callback
	localCallbacks       []schema.Callback
	runID                string
	parentRunID          string
	verbose              bool
}

func newManager(runID string, inheritableCallbacks, localCallbacks []schema.Callback, verbose bool, optFns ...func(*ManagerOptions)) *manager {
	opts := ManagerOptions{}

	for _, fn := range optFns {
		fn(&opts)
	}

	callbacks := append(inheritableCallbacks, localCallbacks...)
	if verbose && !containsWriterCallbackHandler(callbacks) {
		callbacks = append(callbacks, NewWriterHandler())
	}

	return &manager{
		callbacks:            callbacks,
		inheritableCallbacks: inheritableCallbacks,
		localCallbacks:       localCallbacks,
		runID:                runID,
		parentRunID:          opts.ParentRunID,
		verbose:              verbose,
	}
}

func NewManager(inheritableCallbacks, localCallbacks []schema.Callback, verbose bool, optFns ...func(*ManagerOptions)) schema.CallbackManager {
	return newManager("", inheritableCallbacks, localCallbacks, verbose, optFns...)
}

func NewManagerForModelRun(runID string, inheritableCallbacks, localCallbacks []schema.Callback, verbose bool, optFns ...func(*ManagerOptions)) schema.CallbackManagerForModelRun {
	return newManager(runID, inheritableCallbacks, localCallbacks, verbose, optFns...)
}

func NewManagerForChainRun(runID string, inheritableCallbacks, localCallbacks []schema.Callback, verbose bool, optFns ...func(*ManagerOptions)) schema.CallbackManagerForChainRun {
	return newManager(runID, inheritableCallbacks, localCallbacks, verbose, optFns...)
}

func NewManagerForToolRun(runID string, inheritableCallbacks, localCallbacks []schema.Callback, verbose bool, optFns ...func(*ManagerOptions)) schema.CallbackManagerForToolRun {
	return newManager(runID, inheritableCallbacks, localCallbacks, verbose, optFns...)
}

func NewManagerForRetrieverRun(runID string, inheritableCallbacks, localCallbacks []schema.Callback, verbose bool, optFns ...func(*ManagerOptions)) schema.CallbackManagerForRetrieverRun {
	return newManager(runID, inheritableCallbacks, localCallbacks, verbose, optFns...)
}

func (m *manager) GetInheritableCallbacks() []schema.Callback {
	return m.inheritableCallbacks
}

func (m *manager) RunID() string {
	return m.runID
}

func (m *manager) OnLLMStart(ctx context.Context, input *schema.LLMStartManagerInput) (schema.CallbackManagerForModelRun, error) {
	runID := uuid.New().String()

	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnLLMStart(ctx, &schema.LLMStartInput{
				LLMStartManagerInput: input,
				RunID:                runID,
			}); err != nil {
				if c.RaiseError() {
					return nil, err
				}
			}
		}
	}

	return NewManagerForModelRun(runID, m.inheritableCallbacks, m.localCallbacks, m.verbose), nil
}

func (m *manager) OnChatModelStart(ctx context.Context, input *schema.ChatModelStartManagerInput) (schema.CallbackManagerForModelRun, error) {
	runID := uuid.New().String()

	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnChatModelStart(ctx, &schema.ChatModelStartInput{
				ChatModelStartManagerInput: input,
				RunID:                      runID,
			}); err != nil {
				if c.RaiseError() {
					return nil, err
				}
			}
		}
	}

	return NewManagerForModelRun(runID, m.inheritableCallbacks, m.localCallbacks, m.verbose), nil
}

func (m *manager) OnModelNewToken(ctx context.Context, input *schema.ModelNewTokenManagerInput) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnModelNewToken(ctx, &schema.ModelNewTokenInput{
				ModelNewTokenManagerInput: input,
				RunID:                     m.runID,
			}); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnModelEnd(ctx context.Context, input *schema.ModelEndManagerInput) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnModelEnd(ctx, &schema.ModelEndInput{
				ModelEndManagerInput: input,
				RunID:                m.runID,
			}); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnModelError(ctx context.Context, input *schema.ModelErrorManagerInput) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnModelError(ctx, &schema.ModelErrorInput{
				ModelErrorManagerInput: input,
				RunID:                  m.runID,
			}); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnChainStart(ctx context.Context, input *schema.ChainStartManagerInput) (schema.CallbackManagerForChainRun, error) {
	runID := uuid.New().String()

	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnChainStart(ctx, &schema.ChainStartInput{
				ChainStartManagerInput: input,
				RunID:                  runID,
			}); err != nil {
				if c.RaiseError() {
					return nil, err
				}
			}
		}
	}

	return NewManagerForChainRun(runID, m.inheritableCallbacks, m.localCallbacks, m.verbose), nil
}

func (m *manager) OnChainEnd(ctx context.Context, input *schema.ChainEndManagerInput) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnChainEnd(ctx, &schema.ChainEndInput{
				ChainEndManagerInput: input,
				RunID:                m.runID,
			}); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnChainError(ctx context.Context, input *schema.ChainErrorManagerInput) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnChainError(ctx, &schema.ChainErrorInput{
				ChainErrorManagerInput: input,
				RunID:                  m.runID,
			}); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnAgentAction(ctx context.Context, input *schema.AgentActionManagerInput) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnAgentAction(ctx, &schema.AgentActionInput{
				AgentActionManagerInput: input,
				RunID:                   m.runID,
			}); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnAgentFinish(ctx context.Context, input *schema.AgentFinishManagerInput) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnAgentFinish(ctx, &schema.AgentFinishInput{
				AgentFinishManagerInput: input,
				RunID:                   m.runID,
			}); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnToolStart(ctx context.Context, input *schema.ToolStartManagerInput) (schema.CallbackManagerForToolRun, error) {
	runID := uuid.New().String()

	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnToolStart(ctx, &schema.ToolStartInput{
				ToolStartManagerInput: input,
				RunID:                 runID,
			}); err != nil {
				if c.RaiseError() {
					return nil, err
				}
			}
		}
	}

	return NewManagerForToolRun(runID, m.inheritableCallbacks, m.localCallbacks, m.verbose), nil
}

func (m *manager) OnToolEnd(ctx context.Context, input *schema.ToolEndManagerInput) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnToolEnd(ctx, &schema.ToolEndInput{
				ToolEndManagerInput: input,
				RunID:               m.runID,
			}); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnToolError(ctx context.Context, input *schema.ToolErrorManagerInput) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnToolError(ctx, &schema.ToolErrorInput{
				ToolErrorManagerInput: input,
				RunID:                 m.runID,
			}); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnText(ctx context.Context, input *schema.TextManagerInput) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnText(ctx, &schema.TextInput{
				TextManagerInput: input,
				RunID:            m.runID,
			}); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnRetrieverStart(ctx context.Context, input *schema.RetrieverStartManagerInput) (schema.CallbackManagerForRetrieverRun, error) {
	runID := uuid.New().String()

	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnRetrieverStart(ctx, &schema.RetrieverStartInput{
				RetrieverStartManagerInput: input,
				RunID:                      runID,
			}); err != nil {
				if c.RaiseError() {
					return nil, err
				}
			}
		}
	}

	return NewManagerForRetrieverRun(runID, m.inheritableCallbacks, m.localCallbacks, m.verbose), nil
}

func (m *manager) OnRetrieverEnd(ctx context.Context, input *schema.RetrieverEndManagerInput) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnRetrieverEnd(ctx, &schema.RetrieverEndInput{
				RetrieverEndManagerInput: input,
				RunID:                    m.runID,
			}); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnRetrieverError(ctx context.Context, input *schema.RetrieverErrorManagerInput) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnRetrieverError(ctx, &schema.RetrieverErrorInput{
				RetrieverErrorManagerInput: input,
				RunID:                      m.runID,
			}); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func containsWriterCallbackHandler(handlers []schema.Callback) bool {
	for _, handler := range handlers {
		if _, ok := handler.(*WriterHandler); ok {
			return true
		}
	}

	return false
}

// Compile time check to ensure NoopManager satisfies the CallbackManagerForChainRun interface.
var _ schema.CallbackManagerForChainRun = (*NoopManager)(nil)

// Compile time check to ensure NoopManager satisfies the CallbackManagerForModelRun interface.
var _ schema.CallbackManagerForModelRun = (*NoopManager)(nil)

// Compile time check to ensure NoopManager satisfies the CallbackManagerForToolRun interface.
var _ schema.CallbackManagerForToolRun = (*NoopManager)(nil)

type NoopManager struct{}

func (m *NoopManager) OnChainEnd(ctx context.Context, input *schema.ChainEndManagerInput) error {
	return nil
}
func (m *NoopManager) OnChainError(ctx context.Context, input *schema.ChainErrorManagerInput) error {
	return nil
}
func (m *NoopManager) OnAgentAction(ctx context.Context, input *schema.AgentActionManagerInput) error {
	return nil
}
func (m *NoopManager) OnAgentFinish(ctx context.Context, input *schema.AgentFinishManagerInput) error {
	return nil
}

func (m *NoopManager) OnModelNewToken(ctx context.Context, input *schema.ModelNewTokenManagerInput) error {
	return nil
}

func (m *NoopManager) OnModelEnd(ctx context.Context, input *schema.ModelEndManagerInput) error {
	return nil
}

func (m *NoopManager) OnModelError(ctx context.Context, input *schema.ModelErrorManagerInput) error {
	return nil
}

func (m *NoopManager) OnToolEnd(ctx context.Context, input *schema.ToolEndManagerInput) error {
	return nil
}

func (m *NoopManager) OnToolError(ctx context.Context, input *schema.ToolErrorManagerInput) error {
	return nil
}

func (m *NoopManager) OnText(ctx context.Context, input *schema.TextManagerInput) error {
	return nil
}

func (m *NoopManager) GetInheritableCallbacks() []schema.Callback {
	return nil
}

func (m *NoopManager) RunID() string {
	return ""
}
