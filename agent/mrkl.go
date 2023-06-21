package agent

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure ZeroShotReactDescription satisfies the agent interface.
var _ schema.Agent = (*ZeroShotReactDescription)(nil)

const (
	defaultMRKLPrefix = `Answer the following questions as best you can. You have access to the following tools:
{{.toolDescriptions}}`

	defaultMRKLInstructions = `Use the following format:

Question: the input question you must answer
Thought: you should always think about what to do
Action: the action to take, should be one of [{{.toolNames}}]
Action Input: the input to the action
Observation: the result of the action
... (this Thought/Action/Action Input/Observation can repeat N times)
Thought: I now know the final answer
Final Answer: the final answer to the original input question`

	defaultMRKLSuffix = `Begin!

Question: {{.input}}
Thought: {{.agentScratchpad}}`

	finalAnswerAction = "Final Answer:"
)

type ZeroShotReactDescriptionOptions struct {
	Prefix       string
	Instructions string
	Suffix       string
	OutputKey    string
}

type ZeroShotReactDescription struct {
	chain schema.Chain
	tools []schema.Tool
	opts  ZeroShotReactDescriptionOptions
}

func NewZeroShotReactDescription(llm schema.LLM, tools []schema.Tool) (*ZeroShotReactDescription, error) {
	opts := ZeroShotReactDescriptionOptions{
		Prefix:       defaultMRKLPrefix,
		Instructions: defaultMRKLInstructions,
		Suffix:       defaultMRKLSuffix,
		OutputKey:    "output",
	}

	prompt, err := createMRKLPrompt(tools, opts.Prefix, opts.Instructions, opts.Suffix)
	if err != nil {
		return nil, err
	}

	llmChain, err := chain.NewLLMChain(llm, prompt)
	if err != nil {
		return nil, err
	}

	return &ZeroShotReactDescription{
		chain: llmChain,
		tools: tools,
		opts:  opts,
	}, nil
}

func (a *ZeroShotReactDescription) Plan(ctx context.Context, intermediateSteps []schema.AgentStep, inputs map[string]string) ([]schema.AgentAction, *schema.AgentFinish, error) {
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

func (a *ZeroShotReactDescription) InputKeys() []string {
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

func (a *ZeroShotReactDescription) OutputKeys() []string {
	return []string{a.opts.OutputKey}
}

// constructScratchPad constructs the scratchpad that lets the agent
// continue its thought process.
func (a *ZeroShotReactDescription) constructScratchPad(steps []schema.AgentStep) string {
	scratchPad := ""
	for _, step := range steps {
		scratchPad += step.Action.Log
		scratchPad += fmt.Sprintf("\nObservation: %s\nThought:", step.Observation)
	}

	return scratchPad
}

func (a *ZeroShotReactDescription) parseOutput(output string) ([]schema.AgentAction, *schema.AgentFinish, error) {
	if strings.Contains(output, finalAnswerAction) {
		splits := strings.Split(output, finalAnswerAction)

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

func createMRKLPrompt(tools []schema.Tool, prefix, instructions, suffix string) (*prompt.Template, error) {
	return prompt.NewTemplate(strings.Join([]string{prefix, instructions, suffix}, "\n\n"), func(o *prompt.TemplateOptions) {
		o.PartialValues = prompt.PartialValues{
			"toolNames":        toolNames(tools),
			"toolDescriptions": toolDescriptions(tools),
		}
	})
}
