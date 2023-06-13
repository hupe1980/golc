package agent

import (
	"context"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/prompt"
)

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
)

type ZeroShotReactDescriptionAgentOptions struct {
	Prefix       string
	Instructions string
	Suffix       string
}

type ZeroShotReactDescriptionAgent struct {
	chain golc.Chain
	tools []golc.Tool
}

func NewZeroShotReactDescriptionAgent(llm golc.LLM, tools []golc.Tool) (*ZeroShotReactDescriptionAgent, error) {
	opts := ZeroShotReactDescriptionAgentOptions{
		Prefix:       defaultMRKLPrefix,
		Instructions: defaultMRKLInstructions,
		Suffix:       defaultMRKLSuffix,
	}

	prompt, err := createMRKLPrompt(tools, opts.Prefix, opts.Instructions, opts.Suffix)
	if err != nil {
		return nil, err
	}

	llmChain, err := chain.NewLLMChain(llm, prompt)
	if err != nil {
		return nil, err
	}

	return &ZeroShotReactDescriptionAgent{
		chain: llmChain,
		tools: tools,
	}, nil
}

func (a *ZeroShotReactDescriptionAgent) Plan(ctx context.Context) {}

func (a *ZeroShotReactDescriptionAgent) InputKeys() []string {
	return nil
}

func (a *ZeroShotReactDescriptionAgent) OutputKeys() []string {
	return nil
}

func createMRKLPrompt(tools []golc.Tool, prefix, instructions, suffix string) (*prompt.Template, error) {
	return prompt.NewTemplate(strings.Join([]string{prefix, instructions, suffix}, "\n\n"), func(o *prompt.TemplateOptions) {
		o.PartialValues = prompt.PartialValues{
			"toolNames":        toolNames(tools),
			"toolDescriptions": toolDescriptions(tools),
		}
	})
}
