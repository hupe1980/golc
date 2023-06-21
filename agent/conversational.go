package agent

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/memory"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

const (
	defaultConversationalPrefix = `Assistant is a large language model trained by OpenAI.

Assistant is designed to be able to assist with a wide range of tasks, from answering simple questions to providing in-depth explanations and discussions on a wide range of topics. As a language model, Assistant is able to generate human-like text based on the input it receives, allowing it to engage in natural-sounding conversations and provide responses that are coherent and relevant to the topic at hand.

Assistant is constantly learning and improving, and its capabilities are constantly evolving. It is able to process and understand large amounts of text, and can use this knowledge to provide accurate and informative responses to a wide range of questions. Additionally, Assistant is able to generate its own text based on the input it receives, allowing it to engage in discussions and provide explanations and descriptions on a wide range of topics.

Overall, Assistant is a powerful tool that can help with a wide range of tasks and provide valuable insights and information on a wide range of topics. Whether you need help with a specific question or just want to have a conversation about a particular topic, Assistant is here to assist.

TOOLS:
------

Assistant has access to the following tools:`
	defaultConversationalInstructions = `To use a tool, please use the following format:

Thought: Do I need to use a tool? Yes
Action: the action to take, should be one of [{{.toolNames}}]
Action Input: the input to the action
Observation: the result of the action

When you have a response to say to the Human, or if you do not need to use a tool, you MUST use the format:

Thought: Do I need to use a tool? No
{{.aiPrefix}}: [your response here]`

	defaultConversationalSuffix = `Begin!

Previous conversation history:
{{.chatHistory}}

New input: {{.input}}

Thought:{{.agentScratchpad}}`
)

type ConversationalReactDescriptionOptions struct {
	Prefix       string
	Instructions string
	Suffix       string
	AIPrefix     string
	OutputKey    string
}

type ConversationalReactDescription struct {
	chain schema.Chain
	tools []schema.Tool
	opts  ConversationalReactDescriptionOptions
}

func NewConversationalReactDescription(llm schema.LLM, tools []schema.Tool) (*ConversationalReactDescription, error) {
	opts := ConversationalReactDescriptionOptions{
		Prefix:       defaultConversationalPrefix,
		Instructions: defaultConversationalInstructions,
		Suffix:       defaultConversationalSuffix,
		AIPrefix:     "AI",
	}

	prompt, err := createConversationalPrompt(tools, opts.Prefix, opts.Instructions, opts.Suffix)
	if err != nil {
		return nil, err
	}

	llmChain, err := chain.NewLLMChain(llm, prompt, func(o *chain.LLMChainOptions) {
		o.Memory = memory.NewConversationBuffer()
	})
	if err != nil {
		return nil, err
	}

	return &ConversationalReactDescription{
		chain: llmChain,
		tools: tools,
		opts:  opts,
	}, nil
}

func (a *ConversationalReactDescription) Plan(ctx context.Context, intermediateSteps []schema.AgentStep, inputs map[string]string) ([]schema.AgentAction, *schema.AgentFinish, error) {
	fullInputes := make(schema.ChainValues, len(inputs))
	for key, value := range inputs {
		fullInputes[key] = value
	}

	fullInputes["agentScratchpad"] = a.constructScratchPad(intermediateSteps)

	resp, err := golc.Call(ctx, a.chain, fullInputes)
	if err != nil {
		return nil, nil, err
	}

	output, ok := resp[a.chain.OutputKeys()[0]].(string)
	if !ok {
		return nil, nil, ErrInvalidChainReturnType
	}

	return a.parseOutput(output)
}

func (a *ConversationalReactDescription) InputKeys() []string {
	chainInputs := a.chain.InputKeys()

	agentInput := make([]string, 0, len(chainInputs))

	for _, v := range chainInputs {
		if v == "agentScratchpad" {
			continue
		}

		agentInput = append(agentInput, v)
	}

	return agentInput
}

func (a *ConversationalReactDescription) OutputKeys() []string {
	return []string{a.opts.OutputKey}
}

// constructScratchPad constructs the scratchpad that lets the agent
// continue its thought process.
func (a *ConversationalReactDescription) constructScratchPad(steps []schema.AgentStep) string {
	scratchPad := ""
	for _, step := range steps {
		scratchPad += step.Action.Log
		scratchPad += fmt.Sprintf("\nObservation: %s\nThought:", step.Observation)
	}

	return scratchPad
}

func (a *ConversationalReactDescription) parseOutput(output string) ([]schema.AgentAction, *schema.AgentFinish, error) {
	if strings.Contains(output, a.opts.AIPrefix) {
		splits := strings.Split(output, a.opts.AIPrefix)

		return nil, &schema.AgentFinish{
			ReturnValues: map[string]any{
				a.opts.OutputKey: splits[len(splits)-1],
			},
			Log: output,
		}, nil
	}

	r := regexp.MustCompile(`Action:\s*(.+)\s*Action Input:\s*(.+)`)
	matches := r.FindStringSubmatch(output)

	if len(matches) == 0 {
		return nil, nil, fmt.Errorf("%w: %s", ErrUnableToParseOutput, output)
	}

	return []schema.AgentAction{
		{Tool: strings.TrimSpace(matches[1]), ToolInput: strings.TrimSpace(matches[2]), Log: output},
	}, nil, nil
}

func createConversationalPrompt(tools []schema.Tool, prefix, instructions, suffix string) (*prompt.Template, error) {
	return prompt.NewTemplate(strings.Join([]string{prefix, instructions, suffix}, "\n\n"), func(o *prompt.TemplateOptions) {
		o.PartialValues = prompt.PartialValues{
			"toolNames":        toolNames(tools),
			"toolDescriptions": toolDescriptions(tools),
			"chatHistory":      "",
		}
	})
}
