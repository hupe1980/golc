package prompt

import (
	"errors"
	"fmt"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/util"
)

var (
	ErrInvalidPartialVariableType = errors.New("invalid partial variable type")
)

type StringPromptValue string

func (v StringPromptValue) String() string {
	return string(v)
}

func (v StringPromptValue) Messages() []golc.ChatMessage {
	return []golc.ChatMessage{
		golc.NewHumanChatMessage(string(v)),
	}
}

type PartialValues map[string]any

type TemplateOptions struct {
	PartialValues PartialValues
	Language      string
}

type Template struct {
	template      string
	partialValues PartialValues
	language      string
	formatter     *Formatter
}

func NewTemplate(template string, optFns ...func(o *TemplateOptions)) (*Template, error) {
	opts := TemplateOptions{
		Language: "en",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	f, err := NewFormatter(template)
	if err != nil {
		return nil, err
	}

	return &Template{
		template:      template,
		partialValues: opts.PartialValues,
		language:      opts.Language,
		formatter:     f,
	}, nil
}

func (p *Template) Partial(values PartialValues) (*Template, error) {
	return NewTemplate(p.template, func(o *TemplateOptions) {
		o.Language = p.language
		o.PartialValues = util.MergeMaps(p.partialValues, values)
	})
}

func (p *Template) Format(values map[string]any) (string, error) {
	resolvedValues, err := p.resolvePartialValues()
	if err != nil {
		return "", err
	}

	return p.formatter.Render(util.MergeMaps(resolvedValues, values))
}

func (p *Template) resolvePartialValues() (map[string]any, error) {
	resolvedValues := make(map[string]any)

	for variable, value := range p.partialValues {
		switch value := value.(type) {
		case string:
			resolvedValues[variable] = value
		case func() string:
			resolvedValues[variable] = value()
		default:
			return nil, fmt.Errorf("%w: %v", ErrInvalidPartialVariableType, variable)
		}
	}

	return resolvedValues, nil
}

func (p *Template) FormatPrompt(values map[string]any) (golc.PromptValue, error) {
	prompt, err := p.Format(values)
	if err != nil {
		return nil, err
	}

	return StringPromptValue(prompt), nil
}
