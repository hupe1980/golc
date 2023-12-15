package prompt

import (
	"fmt"
	"regexp"

	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Template satisfies the PromptTemplate interface.
var _ schema.PromptTemplate = (*Template)(nil)

// TemplateOptions defines the options for configuring a Template.
type TemplateOptions struct {
	PartialValues           map[string]any
	Language                string
	OutputParser            schema.OutputParser[any]
	TransformPythonTemplate bool
	FormatterOptions
}

var DefaultTemplateOptions = TemplateOptions{
	Language:                "en",
	TransformPythonTemplate: false,
	FormatterOptions: FormatterOptions{
		IgnoreMissingKeys: false,
	},
}

// Template represents a template that can be formatted with dynamic values.
type Template struct {
	template  string
	formatter *Formatter
	opts      TemplateOptions
}

// NewTemplate creates a new Template with the provided template and options.
func NewTemplate(template string, optFns ...func(o *TemplateOptions)) *Template {
	opts := DefaultTemplateOptions

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.TransformPythonTemplate {
		re := regexp.MustCompile(`{([^{}]+)}`)
		template = re.ReplaceAllString(template, "{{.$1}}")
	}

	return &Template{
		template: template,
		formatter: NewFormatter(template, func(o *FormatterOptions) {
			o.IgnoreMissingKeys = opts.IgnoreMissingKeys
			o.TemplateFuncMap = opts.TemplateFuncMap
		}),
		opts: opts,
	}
}

// Partial creates a new Template with partial values.
func (p *Template) Partial(values map[string]any) schema.PromptTemplate {
	return NewTemplate(p.template, func(o *TemplateOptions) {
		o.Language = p.opts.Language
		o.OutputParser = p.opts.OutputParser
		o.PartialValues = util.MergeMaps(p.opts.PartialValues, values)
	})
}

// Format applies values to the template and returns the formatted result.
func (p *Template) Format(values map[string]any) (string, error) {
	resolvedValues, err := p.resolvePartialValues()
	if err != nil {
		return "", err
	}

	return p.formatter.Render(util.MergeMaps(resolvedValues, values))
}

// OutputParser returns the output parser function and a boolean indicating if an output parser is defined.
func (p *Template) OutputParser() (schema.OutputParser[any], bool) {
	if p.opts.OutputParser != nil {
		return p.opts.OutputParser, true
	}

	return nil, false
}

// InputVariables returns the input variables used in the template.
func (p *Template) InputVariables() []string {
	fields := p.formatter.Fields()

	vars := []string{}

	for _, f := range fields {
		name := extractNameFromField(f)
		if name != "" {
			if _, ok := p.opts.PartialValues[name]; !ok {
				vars = append(vars, name)
			}
		}
	}

	return vars
}

// resolvePartialValues resolves partial values to be used in the template.
func (p *Template) resolvePartialValues() (map[string]any, error) {
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

// FormatPrompt applies values to the template and returns a PromptValue representation of the formatted result.
func (p *Template) FormatPrompt(values map[string]any) (schema.PromptValue, error) {
	prompt, err := p.Format(values)
	if err != nil {
		return nil, err
	}

	return StringPromptValue(prompt), nil
}
