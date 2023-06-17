package callback

import (
	"github.com/google/uuid"
	"github.com/hupe1980/golc"
)

type Manager struct {
	callbacks []Callback
	runID     uuid.UUID
	verbose   bool
}

func NewManager(callbacks []Callback, verbose bool) *Manager {
	return &Manager{
		callbacks: callbacks,
		runID:     uuid.New(),
		verbose:   verbose,
	}
}

func (m *Manager) OnLLMStart(llmName string, prompts []string) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnLLMStart(llmName, prompts); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *Manager) OnLLMNewToken(token string) error {
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

func (m *Manager) OnLLMEnd(result golc.LLMResult) error {
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

func (m *Manager) OnLLMError(llmError error) error {
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

func (m *Manager) OnChainStart(chainName string, inputs *golc.ChainValues) error {
	for _, c := range m.callbacks {
		if m.verbose || c.AlwaysVerbose() {
			if err := c.OnChainStart(chainName, inputs); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *Manager) OnChainEnd(outputs *golc.ChainValues) error {
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
