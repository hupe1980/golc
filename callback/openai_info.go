package callback

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/hupe1980/golc/schema"
)

var modelCostMapping = map[string]float64{
	"gpt-4":                     0.03,
	"gpt-4-0314":                0.03,
	"gpt-4-completion":          0.06,
	"gpt-4-0314-completion":     0.06,
	"gpt-4-32k":                 0.06,
	"gpt-4-32k-0314":            0.06,
	"gpt-4-32k-completion":      0.12,
	"gpt-4-32k-0314-completion": 0.12,
	"gpt-3.5-turbo":             0.002,
	"gpt-3.5-turbo-0301":        0.002,
	"text-ada-001":              0.0004,
	"ada":                       0.0004,
	"text-babbage-001":          0.0005,
	"babbage":                   0.0005,
	"text-curie-001":            0.002,
	"curie":                     0.002,
	"text-davinci-003":          0.02,
	"text-davinci-002":          0.02,
	"code-davinci-002":          0.02,
	"ada-finetuned":             0.0016,
	"babbage-finetuned":         0.0024,
	"curie-finetuned":           0.012,
	"davinci-finetuned":         0.12,
}

// Compile time check to ensure OpenAIHandler satisfies the Callback interface.
var _ schema.Callback = (*OpenAIHandler)(nil)

type OpenAIHandler struct {
	NoopHandler
	totalTokens        int
	promptTokens       int
	completionTokens   int
	successfulRequests int
	totalCost          float64
	mu                 sync.Mutex
}

func NewOpenAIHandler() *OpenAIHandler {
	return &OpenAIHandler{}
}

func (cb *OpenAIHandler) String() string {
	return fmt.Sprintf("Tokens Used: %d\nPrompt Tokens: %d\nCompletion Tokens: %d\nSuccessful Requests: %d\nTotal Cost (USD): $%.2f",
		cb.totalTokens, cb.promptTokens, cb.completionTokens, cb.successfulRequests, cb.totalCost)
}

func (cb *OpenAIHandler) AlwaysVerbose() bool {
	return true
}

func (cb *OpenAIHandler) OnModelEnd(ctx context.Context, input *schema.ModelEndInput) error {
	if input.Result.LLMOutput == nil {
		return nil
	}

	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.successfulRequests++

	tokenUsage, ok := input.Result.LLMOutput["TokenUsage"].(map[string]int)
	if !ok {
		return nil
	}

	totalTokens := tokenUsage["TotalTokens"]
	promptTokens := tokenUsage["PromptTokens"]
	completionTokens := tokenUsage["CompletionTokens"]

	if modelName, ok := input.Result.LLMOutput["modelName"].(string); ok {
		completionCosts, err := calculateOpenAITokenCostForModel(modelName, completionTokens, true)
		if err != nil {
			return err
		}

		promptCosts, err := calculateOpenAITokenCostForModel(modelName, promptTokens, false)
		if err != nil {
			return err
		}

		cb.totalCost += completionCosts + promptCosts
	}

	cb.totalTokens += totalTokens
	cb.promptTokens += promptTokens
	cb.completionTokens += completionTokens

	return nil
}

func calculateOpenAITokenCostForModel(modelName string, numTokens int, isCompletion bool) (float64, error) {
	modelName = standardizeModelName(modelName, isCompletion)

	costPer1KTokens, ok := modelCostMapping[modelName]
	if !ok {
		return 0, fmt.Errorf("unknown model: %s", modelName)
	}

	return costPer1KTokens * float64(numTokens) / 1000, nil
}

func standardizeModelName(modelName string, isCompletion bool) string {
	modelName = strings.ToLower(modelName)

	if strings.Contains(modelName, "ft-") {
		return strings.Split(modelName, ":")[0] + "-finetuned"
	} else if isCompletion && strings.HasPrefix(modelName, "gpt-4") {
		return modelName + "-completion"
	} else {
		return modelName
	}
}
