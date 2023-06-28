package chatmodel

import (
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
	"github.com/sashabaranov/go-openai"
)

type AzureOpenAIOptions struct {
	OpenAIOptions
	Deployment string
}

func NewAzureOpenAI(apiKey, baseURL string, optFns ...func(o *AzureOpenAIOptions)) (*OpenAI, error) {
	opts := AzureOpenAIOptions{
		OpenAIOptions: OpenAIOptions{
			CallbackOptions: &schema.CallbackOptions{
				Verbose: golc.Verbose,
			},
			ModelName:        openai.GPT3Dot5Turbo,
			Temperatur:       1,
			TopP:             1,
			PresencePenalty:  0,
			FrequencyPenalty: 0,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Tokenizer == nil {
		opts.Tokenizer = tokenizer.NewOpenAI(opts.ModelName)
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

	return newOpenAI(openai.NewClientWithConfig(config), opts.OpenAIOptions)
}
