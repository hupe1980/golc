package agent

import (
	"context"
	"errors"
	"fmt"

	"github.com/hupe1980/golc/model"
	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tool"
)

// Compile time check to ensure OpenAIFunctions satisfies the agent interface.
var _ schema.Agent = (*OpenAIFunctions)(nil)

type OpenAIFunctionsOptions struct {
	OutputKey string
}

type OpenAIFunctions struct {
	model     schema.ChatModel
	functions []schema.FunctionDefinition
	opts      OpenAIFunctionsOptions
}

func NewOpenAIFunctions(model schema.Model, tools []schema.Tool) (*Executor, error) {
	opts := OpenAIFunctionsOptions{
		OutputKey: "output",
	}

	chatModel, ok := model.(*chatmodel.OpenAI)
	if !ok {
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

	agent := &OpenAIFunctions{
		model:     chatModel,
		functions: functions,
		opts:      opts,
	}

	return NewExecutor(agent, tools)
}

func (a *OpenAIFunctions) Plan(ctx context.Context, intermediateSteps []schema.AgentStep, inputs schema.ChainValues) ([]*schema.AgentAction, *schema.AgentFinish, error) {
	inputs["agentScratchpad"] = a.constructScratchPad(intermediateSteps)

	chatTemplate := prompt.NewChatTemplate([]prompt.MessageTemplate{
		prompt.NewSystemMessageTemplate("You are a helpful AI assistant."),
		prompt.NewHumanMessageTemplate("{{.input}}"),
	})

	placeholder := prompt.NewMessagesPlaceholder("agentScratchpad")

	wrapper := prompt.NewChatTemplateWrapper(chatTemplate, placeholder)

	prompt, err := wrapper.FormatPrompt(inputs)
	if err != nil {
		return nil, nil, err
	}

	result, err := model.ChatModelGenerate(ctx, a.model, prompt.Messages(), func(o *model.Options) {
		o.Functions = a.functions
	})
	if err != nil {
		return nil, nil, err
	}

	msg := result.Generations[0].Message

	aiMsg, ok := msg.(*schema.AIChatMessage)
	if !ok {
		return nil, nil, fmt.Errorf("unexpected chatMessage type. Expected ai, but got %s", msg.Type())
	}

	ext := aiMsg.Extension()

	if ext.FunctionCall != nil {
		toolInput := schema.NewToolInputFromArguments(ext.FunctionCall.Arguments)

		msgContent := ""
		if aiMsg.Content() != "" {
			msgContent = fmt.Sprintf("responded: %s", aiMsg.Content())
		}

		log := fmt.Sprintf("\nInvoking `%s` with `%s`\n%s\n", ext.FunctionCall.Name, toolInput, msgContent)

		return []*schema.AgentAction{
			{Tool: ext.FunctionCall.Name, ToolInput: toolInput, Log: log, MessageLog: schema.ChatMessages{aiMsg}},
		}, nil, nil
	}

	return nil, &schema.AgentFinish{
		ReturnValues: map[string]any{
			a.opts.OutputKey: aiMsg.Content(),
		},
		Log: aiMsg.Content(),
	}, nil
}

func (a *OpenAIFunctions) InputKeys() []string {
	return []string{"input"}
}

func (a *OpenAIFunctions) OutputKeys() []string {
	return []string{a.opts.OutputKey}
}

func (a *OpenAIFunctions) constructScratchPad(steps []schema.AgentStep) schema.ChatMessages {
	messages := schema.ChatMessages{}

	for _, step := range steps {
		if step.Action.MessageLog != nil {
			messages = append(messages, step.Action.MessageLog...)
			messages = append(messages, schema.NewFunctionChatMessage(step.Action.Tool, step.Observation))
		} else {
			messages = append(messages, schema.NewAIChatMessage(step.Action.Log))
		}
	}

	return messages
}
