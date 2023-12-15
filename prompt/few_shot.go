package prompt

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure FewShotTemplate satisfies the PromptTemplate interface.
var _ schema.PromptTemplate = (*FewShotTemplate)(nil)

// FewShotTemplateOptions represents options for configuring a FewShotTemplate.
type FewShotTemplateOptions struct {
	// Prefix to be added before the template.
	Prefix string
	// Separator between examples and the template.
	Separator string
	// OutputParser to parse the response.
	OutputParser schema.OutputParser[any]
	// PartialValues to be used in the template.
	PartialValues map[string]any
	// IgnoreMissingKeys allows ignoring missing keys in the template.
	IgnoreMissingKeys bool
}

// FewShotTemplate is a template that combines examples with a main template.
type FewShotTemplate struct {
	template        string
	examples        []map[string]any
	exampleTemplate *Template
	opts            FewShotTemplateOptions
}

// NewFewShotTemplate creates a new FewShotTemplate with the provided template, examples, and options.
func NewFewShotTemplate(template string, examples []map[string]any, exampleTemplate *Template, optFns ...func(o *FewShotTemplateOptions)) *FewShotTemplate {
	opts := FewShotTemplateOptions{
		Separator:         "\n\n",
		IgnoreMissingKeys: false,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &FewShotTemplate{
		template:        template,
		examples:        examples,
		exampleTemplate: exampleTemplate,
		opts:            opts,
	}
}

// Format applies values to the template and returns the formatted result.
func (p *FewShotTemplate) Format(values map[string]any) (string, error) {
	pieces := []string{}

	if p.opts.Prefix != "" {
		pieces = append(pieces, p.opts.Prefix)
	}

	for _, example := range p.examples {
		e, err := p.exampleTemplate.Format(example)
		if err != nil {
			return "", err
		}

		pieces = append(pieces, e)
	}

	pieces = append(pieces, p.template)

	formatter := NewFormatter(strings.Join(pieces, p.opts.Separator), func(o *FormatterOptions) {
		o.IgnoreMissingKeys = p.opts.IgnoreMissingKeys
	})

	resolvedValues, err := p.resolvePartialValues()
	if err != nil {
		return "", err
	}

	return formatter.Render(util.MergeMaps(resolvedValues, values))
}

// FormatPrompt applies values to the template and returns a PromptValue representation of the formatted result.
func (p *FewShotTemplate) FormatPrompt(values map[string]any) (schema.PromptValue, error) {
	prompt, err := p.Format(values)
	if err != nil {
		return nil, err
	}

	return StringPromptValue(prompt), nil
}

// Partial creates a new FewShotTemplate with partial values.
func (p *FewShotTemplate) Partial(values map[string]any) schema.PromptTemplate {
	return NewFewShotTemplate(p.template, p.examples, p.exampleTemplate, func(o *FewShotTemplateOptions) {
		o.Prefix = p.opts.Prefix
		o.Separator = p.opts.Separator
		o.OutputParser = p.opts.OutputParser
		o.PartialValues = util.MergeMaps(p.opts.PartialValues, values)
		o.IgnoreMissingKeys = p.opts.IgnoreMissingKeys
	})
}

// OutputParser returns the output parser function and a boolean indicating if an output parser is defined.
func (p *FewShotTemplate) OutputParser() (schema.OutputParser[any], bool) {
	if p.opts.OutputParser != nil {
		return p.opts.OutputParser, true
	}

	return nil, false
}

// InputVariables returns the input variables used in the template.
func (p *FewShotTemplate) InputVariables() []string {
	vars := p.exampleTemplate.InputVariables()

	t := template.Must(template.New("template").Parse(p.template))

	for _, f := range ListTemplateFields(t) {
		name := extractNameFromField(f)
		if name != "" {
			if _, ok := p.opts.PartialValues[name]; !ok {
				if !util.Contains(vars, name) {
					vars = append(vars, name)
				}
			}
		}
	}

	return vars
}

// resolvePartialValues resolves partial values to be used in the template.
func (p *FewShotTemplate) resolvePartialValues() (map[string]any, error) {
	resolvedValues := make(map[string]any)

	for variable, value := range p.opts.PartialValues {
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
