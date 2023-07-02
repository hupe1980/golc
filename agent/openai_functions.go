package agent

import (
	"context"
	"errors"
	"fmt"

	"github.com/hupe1980/golc/model"
	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tool"
)

// Compile time check to ensure OpenAIFunctions satisfies the agent interface.
var _ schema.Agent = (*OpenAIFunctions)(nil)

type OpenAIFunctionsOptions struct{}

type OpenAIFunctions struct {
	model     schema.ChatModel
	functions []schema.FunctionDefinition
	opts      OpenAIFunctionsOptions
}

func NewOpenAIFunctions(model schema.ChatModel, tools []schema.Tool) (*OpenAIFunctions, error) {
	opts := OpenAIFunctionsOptions{}

	if _, ok := model.(*chatmodel.OpenAI); !ok {
		return nil, errors.New("agent only supports OpenAI chatModels")
	}

	functions := make([]schema.FunctionDefinition, len(tools))

	for i, t := range tools {
		f, err := tool.ToFunction(t)
		if err != nil {
			return nil, err
		}

		functions[i] = *f
	}

	return &OpenAIFunctions{
		model:     model,
		functions: functions,
		opts:      opts,
	}, nil
}

func (a *OpenAIFunctions) Plan(ctx context.Context, intermediateSteps []schema.AgentStep, inputs map[string]string) ([]*schema.AgentAction, *schema.AgentFinish, error) {
	fullInputes := make(schema.ChainValues, len(inputs))
	for key, value := range inputs {
		fullInputes[key] = value
	}

	fullInputes["agentScratchpad"] = a.constructScratchPad(intermediateSteps)

	result, err := model.ChatModelGenerate(ctx, a.model, nil, func(o *model.Options) {
		o.Functions = a.functions
	})
	if err != nil {
		return nil, nil, err
	}

	msg := result.Generations[0][0].Message

	aiMsg, ok := msg.(schema.AIChatMessage)
	if !ok {
		return nil, nil, fmt.Errorf("unexpected chatMessage type. Expected ai, but got %s", msg.Type())
	}

	attrs := aiMsg.AdditionalAttributes()

	fmt.Println("TODO", attrs)

	return nil, nil, nil
}

func (a *OpenAIFunctions) InputKeys() []string {
	return []string{"input"}
}

func (a *OpenAIFunctions) OutputKeys() []string {
	return []string{"todo"}
}

func (a *OpenAIFunctions) constructScratchPad(steps []schema.AgentStep) []schema.ChatMessage {
	messages := make([]schema.ChatMessage, len(steps))
	for i, step := range steps {
		messages[i] = schema.NewAIChatMessage(step.Action.Log)
	}

	return messages
}
