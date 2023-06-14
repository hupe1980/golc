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

func (m *Manager) OnLLMStart() *ManagerForLLMRun {
	return NewManagerForLLMRun()
}

type ManagerForLLMRun struct{}

func NewManagerForLLMRun() *ManagerForLLMRun {
	return &ManagerForLLMRun{}
}
