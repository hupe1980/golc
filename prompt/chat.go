package prompt

import "github.com/hupe1980/golc"

type ChatPromptValue struct{}

type ChatTemplate struct{}

type SystemMessageTemplate struct {
	prompt *Template
}

func NewSystemMessageTemplate(prompt *Template) *SystemMessageTemplate {
	return &SystemMessageTemplate{
		prompt: prompt,
	}
}

func (pt *SystemMessageTemplate) Format(values map[string]any) (*golc.SystemChatMessage, error) {
	text, err := pt.prompt.Format(values)
	if err != nil {
		return nil, err
	}

	return golc.NewSystemChatMessage(text), nil
}

type AIMessagetTemplate struct {
	prompt *Template
}

func (pt *AIMessagetTemplate) Format(values map[string]any) (*golc.AIChatMessage, error) {
	text, err := pt.prompt.Format(values)
	if err != nil {
		return nil, err
	}

	return golc.NewAIChatMessage(text), nil
}

type HumanMessageTemplate struct {
	prompt *Template
}

func (pt *HumanMessageTemplate) Format(values map[string]any) (*golc.HumanChatMessage, error) {
	text, err := pt.prompt.Format(values)
	if err != nil {
		return nil, err
	}

	return golc.NewHumanChatMessage(text), nil
}
