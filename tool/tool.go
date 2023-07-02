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

func Run(ctx context.Context, t schema.Tool, input string, optFns ...func(o *Options)) (string, error) {
	opts := Options{}

	for _, fn := range optFns {
		fn(&opts)
	}

	cm := callback.NewManager(opts.Callbacks, nil, false)

	rm, err := cm.OnToolStart(ctx, &schema.ToolStartManagerInput{
		ToolName: t.Name(),
		Input:    input,
	})
	if err != nil {
		return "", err
	}

	output, err := t.Run(ctx, input)
	if err != nil {
		if cbErr := rm.OnToolError(ctx, &schema.ToolErrorManagerInput{
			Error: err,
		}); cbErr != nil {
			return "", cbErr
		}

		return "", err
	}

	if err := rm.OnToolEnd(ctx, &schema.ToolEndManagerInput{
		Output: output,
	}); err != nil {
		return "", err
	}

	return output, nil
}

// ToFunction formats a tool into a function API
func ToFunction(t schema.Tool) (*schema.FunctionDefinition, error) {
	function := &schema.FunctionDefinition{
		Name:        t.Name(),
		Description: t.Description(),
	}

	run := reflect.TypeOf(t.Run)

	in := run.In(1) // ignore context at idx 0
	if in.Kind() == reflect.String {
		function.Parameters = schema.FunctionDefinitionParameters{
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

	jsonSchema, err := jsonschema.Generate(in)
	if err != nil {
		return nil, err
	}

	function.Parameters = schema.FunctionDefinitionParameters{
		Type:       "object",
		Properties: jsonSchema.Properties,
		Required:   jsonSchema.Required,
	}

	return function, nil
}
