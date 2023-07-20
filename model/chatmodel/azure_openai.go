package chatmodel

import (
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/schema"
	"github.com/sashabaranov/go-openai"
)

type AzureOpenAIOptions struct {
	OpenAIOptions
	Deployment string
}

type AzureOpenAI struct {
	*OpenAI
}

func NewAzureOpenAI(apiKey, baseURL string, optFns ...func(o *AzureOpenAIOptions)) (*AzureOpenAI, error) {
	opts := AzureOpenAIOptions{
		OpenAIOptions: OpenAIOptions{
			CallbackOptions: &schema.CallbackOptions{
				Verbose: golc.Verbose,
			},
			ModelName:        openai.GPT3Dot5Turbo,
			Temperature:      1,
			TopP:             1,
			PresencePenalty:  0,
			FrequencyPenalty: 0,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	config := openai.DefaultAzureConfig(apiKey, baseURL)
	if opts.Deployment != "" {
		config.AzureModelMapperFunc = func(model string) string {
			azureModelMapping := map[string]string{
				opts.ModelName: opts.Deployment,
			}

			return azureModelMapping[model]
		}
	}

	openAI, err := NewOpenAIFromClient(openai.NewClientWithConfig(config), func(o *OpenAIOptions) { // nolint staticcheck
		o = &opts.OpenAIOptions // nolint ineffassign
	})
	if err != nil {
		return nil, err
	}

	return &AzureOpenAI{OpenAI: openAI}, nil
}

// Type returns the type of the model.
func (cm *AzureOpenAI) Type() string {
	return "chatmodel.AzureOpenAI"
}
