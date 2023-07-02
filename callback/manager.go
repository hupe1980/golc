package callback

import (
	"github.com/google/uuid"
	"github.com/hupe1980/golc/schema"
)

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
	if verbose && !containsStdOutCallbackHandler(callbacks) {
		callbacks = append(callbacks, NewStdOutHandler())
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

func NewManagerForModelRun(runID string, inheritableCallbacks, localCallbacks []schema.Callback, verbose bool, optFns ...func(*ManagerOptions)) schema.CallBackManagerForModelRun {
	return newManager(runID, inheritableCallbacks, localCallbacks, verbose, optFns...)
}

func NewManagerForChainRun(runID string, inheritableCallbacks, localCallbacks []schema.Callback, verbose bool, optFns ...func(*ManagerOptions)) schema.CallBackManagerForChainRun {
	return newManager(runID, inheritableCallbacks, localCallbacks, verbose, optFns...)
}

func NewManagerForToolRun(runID string, inheritableCallbacks, localCallbacks []schema.Callback, verbose bool, optFns ...func(*ManagerOptions)) schema.CallBackManagerForToolRun {
	return newManager(runID, inheritableCallbacks, localCallbacks, verbose, optFns...)
}

func (m *manager) GetInheritableCallbacks() []schema.Callback {
	return m.inheritableCallbacks
}

func (m *manager) RunID() string {
	return m.runID
}

func (m *manager) OnLLMStart(llmName string, prompts []string, invocationParams map[string]any) (schema.CallBackManagerForModelRun, error) {
	runID := uuid.New().String()

	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnLLMStart(llmName, prompts, invocationParams, runID); err != nil {
				if c.RaiseError() {
					return nil, err
				}
			}
		}
	}

	return NewManagerForModelRun(runID, m.inheritableCallbacks, m.localCallbacks, m.verbose), nil
}

func (m *manager) OnChatModelStart(llmName string, messages []schema.ChatMessages) (schema.CallBackManagerForModelRun, error) {
	runID := uuid.New().String()

	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnChatModelStart(llmName, messages); err != nil {
				if c.RaiseError() {
					return nil, err
				}
			}
		}
	}

	return NewManagerForModelRun(runID, m.inheritableCallbacks, m.localCallbacks, m.verbose), nil
}

func (m *manager) OnModelNewToken(token string) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnModelNewToken(token); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnModelEnd(result schema.ModelResult) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnModelEnd(result, m.runID); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnModelError(llmError error) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnModelError(llmError); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnChainStart(chainName string, inputs schema.ChainValues) (schema.CallBackManagerForChainRun, error) {
	runID := uuid.New().String()

	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnChainStart(chainName, inputs); err != nil {
				if c.RaiseError() {
					return nil, err
				}
			}
		}
	}

	return NewManagerForChainRun(runID, m.inheritableCallbacks, m.localCallbacks, m.verbose), nil
}

func (m *manager) OnChainEnd(outputs schema.ChainValues) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnChainEnd(outputs); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnChainError(chainError error) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnChainError(chainError); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnAgentAction(action schema.AgentAction) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnAgentAction(action); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnAgentFinish(finish schema.AgentFinish) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnAgentFinish(finish); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnToolStart(toolName string, input string) (schema.CallBackManagerForToolRun, error) {
	runID := uuid.New().String()

	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnToolStart(toolName, input); err != nil {
				if c.RaiseError() {
					return nil, err
				}
			}
		}
	}

	return NewManagerForToolRun(runID, m.inheritableCallbacks, m.localCallbacks, m.verbose), nil
}

func (m *manager) OnToolEnd(output string) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnToolEnd(output); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnToolError(toolError error) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnToolError(toolError); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnText(text string) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnText(text); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func containsStdOutCallbackHandler(handlers []schema.Callback) bool {
	for _, handler := range handlers {
		if _, ok := handler.(*StdOutHandler); ok {
			return true
		}
	}

	return false
}
