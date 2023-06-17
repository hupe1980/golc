package prompt

import "github.com/hupe1980/golc/schema"

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

func (pt *SystemMessageTemplate) Format(values map[string]any) (*schema.SystemChatMessage, error) {
	text, err := pt.prompt.Format(values)
	if err != nil {
		return nil, err
	}

	return schema.NewSystemChatMessage(text), nil
}

type AIMessagetTemplate struct {
	prompt *Template
}

func (pt *AIMessagetTemplate) Format(values map[string]any) (*schema.AIChatMessage, error) {
	text, err := pt.prompt.Format(values)
	if err != nil {
		return nil, err
	}

	return schema.NewAIChatMessage(text), nil
}

type HumanMessageTemplate struct {
	prompt *Template
}

func (pt *HumanMessageTemplate) Format(values map[string]any) (*schema.HumanChatMessage, error) {
	text, err := pt.prompt.Format(values)
	if err != nil {
		return nil, err
	}

	return schema.NewHumanChatMessage(text), nil
}
