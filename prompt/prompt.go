package prompt

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/util"
)

var (
	ErrInvalidPartialVariableType = errors.New("invalid partial variable type")
)

type StringPromptValue string

func (v StringPromptValue) String() string {
	return string(v)
}

func (v StringPromptValue) Messages() schema.ChatMessages {
	return schema.ChatMessages{
		schema.NewHumanChatMessage(string(v)),
	}
}

type PartialValues map[string]any

type TemplateOptions struct {
	PartialValues           PartialValues
	Language                string
	OutputParser            schema.OutputParser[any]
	TransformPythonTemplate bool
}

type Template struct {
	template      string
	partialValues PartialValues
	language      string
	outputParser  schema.OutputParser[any]
	formatter     *Formatter
}

func NewTemplate(template string, optFns ...func(o *TemplateOptions)) *Template {
	opts := TemplateOptions{
		Language:                "en",
		TransformPythonTemplate: false,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.TransformPythonTemplate {
		re := regexp.MustCompile(`{([^{}]+)}`)
		template = re.ReplaceAllString(template, "{{.$1}}")
	}

	return &Template{
		template:      template,
		partialValues: opts.PartialValues,
		language:      opts.Language,
		outputParser:  opts.OutputParser,
		formatter:     NewFormatter(template),
	}
}

func (p *Template) Partial(values PartialValues) *Template {
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

func (p *Template) OutputParser() (schema.OutputParser[any], bool) {
	if p.outputParser != nil {
		return p.outputParser, true
	}

	return nil, false
}

func (p *Template) InputVariables() []string {
	fields := p.formatter.Fields()

	vars := []string{}

	for _, f := range fields {
		name := extractNameFromField(f)
		if name != "" {
			if _, ok := p.partialValues[name]; !ok {
				vars = append(vars, name)
			}
		}
	}

	return vars
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

func (p *Template) FormatPrompt(values map[string]any) (schema.PromptValue, error) {
	prompt, err := p.Format(values)
	if err != nil {
		return nil, err
	}

	return StringPromptValue(prompt), nil
}

func extractNameFromField(input string) string {
	re := regexp.MustCompile(`{{\.(.*?)}}`)
	matches := re.FindStringSubmatch(input)

	if len(matches) == 2 {
		return matches[1]
	}

	return ""
}
