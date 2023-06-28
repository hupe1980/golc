package tool

import (
	"context"
	"reflect"

	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/integration/jsonschema"
	"github.com/hupe1980/golc/schema"
)

type Options struct {
	Callbacks   []schema.Callback
	ParentRunID string
}

func Run(ctx context.Context, t schema.Tool, query string, optFns ...func(o *Options)) (string, error) {
	opts := Options{}

	for _, fn := range optFns {
		fn(&opts)
	}

	cm := callback.NewManager(opts.Callbacks, nil, false)

	rm, err := cm.OnToolStart(t.Name(), query)
	if err != nil {
		return "", err
	}

	output, err := t.Run(ctx, query)
	if err != nil {
		if cbErr := rm.OnToolError(err); cbErr != nil {
			return "", cbErr
		}

		return "", err
	}

	if err := rm.OnToolEnd(output); err != nil {
		return "", err
	}

	return output, nil
}

type OpenAIFunctionParameters struct {
	Type       string                        `json:"type"`
	Properties map[string]*jsonschema.Schema `json:"properties"`
	Required   []string                      `json:"required"`
}

type OpenAIFunction struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Parameters  OpenAIFunctionParameters `json:"parameters"`
}

// ToOpenAIFunction formats a tool into the OpenAI function API
func ToOpenAIFunction(t schema.Tool) (*OpenAIFunction, error) {
	function := &OpenAIFunction{
		Name:        t.Name(),
		Description: t.Description(),
	}

	run := reflect.TypeOf(t.Run)

	in := run.In(1) // ignore context at idx 0
	if in.Kind() == reflect.String {
		function.Parameters = OpenAIFunctionParameters{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"__arg1": {
					Type:        "string",
					Description: "__arg1",
				},
			},
			Required: []string{"__arg1"},
		}

		return function, nil
	}

	schema, err := jsonschema.Generate(in)
	if err != nil {
		return nil, err
	}

	function.Parameters = OpenAIFunctionParameters{
		Type:       "object",
		Properties: schema.Properties,
		Required:   schema.Required,
	}

	return function, nil
}
