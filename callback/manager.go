package callback

import (
	"github.com/google/uuid"
	"github.com/hupe1980/golc/schema"
)

type Manager struct {
	callbacks []schema.Callback
	runID     uuid.UUID
	verbose   bool
}

func NewManager(callbacks []schema.Callback, verbose bool) *Manager {
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

func (m *Manager) OnLLMEnd(result *schema.LLMResult) error {
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

func (m *Manager) OnChainStart(chainName string, inputs *schema.ChainValues) error {
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

func (m *Manager) OnChainEnd(outputs *schema.ChainValues) error {
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

func (m *Manager) OnChainError(chainError error) error {
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
