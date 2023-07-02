package callback

import (
	"context"
	"fmt"
	"time"

	"github.com/hupe1980/go-promptlayer"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure  PromptLayerHandler satisfies the Callback interface.
var _ schema.Callback = (*PromptLayerHandler)(nil)

type OnPromptLayerOutputFunc func(output *promptlayer.TrackRequestOutput) error

type PromptLayerHandlerOptions struct {
	PromptID                string
	OnPromptLayerOutputFunc OnPromptLayerOutputFunc
	Tags                    []string
}

type PromptLayerHandler struct {
	handler
	apiKey  string
	client  *promptlayer.Client
	runInfo map[string]map[string]any
	opts    PromptLayerHandlerOptions
}

func NewPromptLayerHandler(apiKey string, optFns ...func(o *PromptLayerHandlerOptions)) *PromptLayerHandler {
	opts := PromptLayerHandlerOptions{}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &PromptLayerHandler{
		apiKey:  apiKey,
		client:  promptlayer.NewClient(apiKey),
		runInfo: map[string]map[string]any{},
		opts:    opts,
	}
}

func (cb PromptLayerHandler) AlwaysVerbose() bool {
	return true
}

func (cb PromptLayerHandler) OnLLMStart(ctx context.Context, input *schema.LLMStartInput) error {
	if cb.opts.PromptID != "" && len(input.Prompts) != 1 {
		panic(fmt.Sprintf("promptID assignment only possible with a single prompt, got %d", len(input.Prompts)))
	}

	if input.LLMType != "llm.OpenAI" {
		panic("currently only openai is supported")
	}

	cb.runInfo[input.RunID] = map[string]any{
		"name":             "openai.Completion.create",
		"prompts":          input.Prompts,
		"invocationParams": input.InvocationParams,
		"startTime":        time.Now(),
	}

	return nil
}

func (cb PromptLayerHandler) OnChatModelStart(ctx context.Context, input *schema.ChatModelStartInput) error {
	return nil
}

func (cb PromptLayerHandler) OnModelEnd(ctx context.Context, input *schema.ModelEndInput) error {
	runInfo, ok := cb.runInfo[input.RunID]
	if !ok {
		return fmt.Errorf("no runInfo for runID %s", input.RunID)
	}

	functionName, _ := runInfo["name"].(string)
	startTime, _ := runInfo["startTime"].(time.Time)
	invocationParams, _ := runInfo["invocationParams"].(map[string]any)
	prompts, _ := runInfo["prompts"].([]string)

	endTime := time.Now()

	for _, generation := range input.Result.Generations {
		output, err := cb.client.TrackRequest(context.Background(), &promptlayer.TrackRequestInput{
			FunctionName: functionName,
			// kwargs will need messages if using chat-based completion
			Kwargs: map[string]any{
				"engine": invocationParams["ModelName"],
				"prompt": prompts[0],
			},
			Tags: cb.opts.Tags,
			RequestResponse: map[string]any{
				"choices": []map[string]any{
					{
						"text": generation[0].Text,
						"info": generation[0].Info,
					},
				},
			},
			PromptID:         cb.opts.PromptID,
			RequestStartTime: startTime,
			RequestEndTime:   endTime,
		})
		if err != nil {
			return err
		}

		if cb.opts.OnPromptLayerOutputFunc != nil {
			if err := cb.opts.OnPromptLayerOutputFunc(output); err != nil {
				return err
			}
		}
	}

	return nil
}
