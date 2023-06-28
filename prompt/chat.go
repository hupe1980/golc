package prompt

import "github.com/hupe1980/golc/schema"

type ChatPromptValue struct{}

type ChatTemplate struct{}

type SystemMessageTemplate struct {
	prompt *Template
}

func NewSystemMessageTemplate(template string) *SystemMessageTemplate {
	return &SystemMessageTemplate{
		prompt: NewTemplate(template),
	}
}

func (pt *SystemMessageTemplate) Format(values map[string]any) (*schema.SystemChatMessage, error) {
	text, err := pt.prompt.Format(values)
	if err != nil {
		return nil, err
	}

	return schema.NewSystemChatMessage(text), nil
}

type AIMessageTemplate struct {
	prompt *Template
}

func NewAIMessageTemplate(template string) *AIMessageTemplate {
	return &AIMessageTemplate{
		prompt: NewTemplate(template),
	}
}

func (pt *AIMessageTemplate) Format(values map[string]any) (*schema.AIChatMessage, error) {
	text, err := pt.prompt.Format(values)
	if err != nil {
		return nil, err
	}

	return schema.NewAIChatMessage(text), nil
}

type HumanMessageTemplate struct {
	prompt *Template
}

func NewHumanMessageTemplate(template string) *HumanMessageTemplate {
	return &HumanMessageTemplate{
		prompt: NewTemplate(template),
	}
}

func (pt *HumanMessageTemplate) Format(values map[string]any) (*schema.HumanChatMessage, error) {
	text, err := pt.prompt.Format(values)
	if err != nil {
		return nil, err
	}

	return schema.NewHumanChatMessage(text), nil
}
