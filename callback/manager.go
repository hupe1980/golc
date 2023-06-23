package callback

import (
	"github.com/google/uuid"
	"github.com/hupe1980/golc/schema"
)

type ManagerOptions struct {
	RunID       string
	ParentRunID string
}

type manager struct {
	callbacks            []schema.Callback
	inheritableCallbacks []schema.Callback
	runID                string
	parentRunID          string
	verbose              bool
}

func newManager(inheritableCallbacks, localCallbacks []schema.Callback, verbose bool, optFns ...func(*ManagerOptions)) *manager {
	opts := ManagerOptions{}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.RunID == "" {
		opts.RunID = uuid.New().String()
	}

	callbacks := append(inheritableCallbacks, localCallbacks...)
	if verbose && !containsStdOutCallbackHandler(callbacks) {
		callbacks = append(callbacks, NewStdOutHandler())
	}

	return &manager{
		callbacks:            callbacks,
		inheritableCallbacks: inheritableCallbacks,
		runID:                opts.RunID,
		parentRunID:          opts.ParentRunID,
		verbose:              verbose,
	}
}

func NewManager(inheritableCallbacks, localCallbacks []schema.Callback, verbose bool, optFns ...func(*ManagerOptions)) schema.CallbackManager {
	return newManager(inheritableCallbacks, localCallbacks, verbose, optFns...)
}

func NewManagerForLLMRun(inheritableCallbacks, localCallbacks []schema.Callback, verbose bool, optFns ...func(*ManagerOptions)) schema.CallBackManagerForLLMRun {
	return newManager(inheritableCallbacks, localCallbacks, verbose, optFns...)
}

func NewManagerForChainRun(inheritableCallbacks, localCallbacks []schema.Callback, verbose bool, optFns ...func(*ManagerOptions)) schema.CallBackManagerForChainRun {
	return newManager(inheritableCallbacks, localCallbacks, verbose, optFns...)
}

func (m *manager) GetInheritableCallbacks() []schema.Callback {
	return m.inheritableCallbacks
}

func (m *manager) RunID() string {
	return m.runID
}

func (m *manager) OnLLMStart(llmName string, prompts []string) (schema.CallBackManagerForLLMRun, error) {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnLLMStart(llmName, prompts); err != nil {
				if c.RaiseError() {
					return nil, err
				}
			}
		}
	}

	return NewManagerForLLMRun(m.inheritableCallbacks, m.callbacks, m.verbose), nil
}

func (m *manager) OnLLMNewToken(token string) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnLLMNewToken(token); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnLLMEnd(result *schema.LLMResult) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnLLMEnd(result); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnLLMError(llmError error) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnLLMError(llmError); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *manager) OnChainStart(chainName string, inputs *schema.ChainValues) (schema.CallBackManagerForChainRun, error) {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnChainStart(chainName, inputs); err != nil {
				if c.RaiseError() {
					return nil, err
				}
			}
		}
	}

	return NewManagerForChainRun(m.inheritableCallbacks, m.callbacks, m.verbose), nil
}

func (m *manager) OnChainEnd(outputs *schema.ChainValues) error {
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

func containsStdOutCallbackHandler(handlers []schema.Callback) bool {
	for _, handler := range handlers {
		if _, ok := handler.(*StdOutHandler); ok {
			return true
		}
	}

	return false
}
