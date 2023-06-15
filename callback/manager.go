package callback

import "github.com/hupe1980/golc"

type Manager struct {
	callbacks []golc.Callback
}

func NewManager(callbacks []golc.Callback) *Manager {
	return &Manager{
		callbacks: callbacks,
	}
}

func (m *Manager) OnLLMStart(verbose bool) error {
	for _, c := range m.callbacks {
		if verbose || c.AlwaysVerbose() {
			if err := c.OnLLMStart(verbose); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *Manager) OnLLMNewToken(token string, verbose bool) error {
	for _, c := range m.callbacks {
		if verbose || c.AlwaysVerbose() {
			if err := c.OnLLMNewToken(token, verbose); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *Manager) OnLLMEnd(result golc.LLMResult, verbose bool) error {
	for _, c := range m.callbacks {
		if verbose || c.AlwaysVerbose() {
			if err := c.OnLLMEnd(result, verbose); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}

func (m *Manager) OnLLMError(e error, verbose bool) error {
	for _, c := range m.callbacks {
		if verbose || c.AlwaysVerbose() {
			if err := c.OnLLMError(e, verbose); err != nil {
				if c.RaiseError() {
					return err
				}
			}
		}
	}

	return nil
}
