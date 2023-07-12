package embedding

import (
	"github.com/sashabaranov/go-openai"
)

type AzureOpenAIOptions struct {
	OpenAIOptions
	APIVersion string
	Deployment string
}

func NewAzureOpenAI(apiKey, baseURL string, optFns ...func(o *AzureOpenAIOptions)) (*OpenAI, error) {
	opts := AzureOpenAIOptions{
		OpenAIOptions: OpenAIOptions{
			ModelName:              "text-embedding-ada-002",
			EmbeddingContextLength: 8191,
			ChunkSize:              1000,
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

	if opts.APIVersion != "" {
		config.APIVersion = opts.APIVersion
	}

	client := openai.NewClientWithConfig(config)

	return NewOpenAIFromClient(client, func(o *OpenAIOptions) { // nolint staticcheck
		o = &opts.OpenAIOptions // nolint ineffassign
	})
}
