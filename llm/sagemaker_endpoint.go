package llm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntime"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/tokenizer"
	"github.com/hupe1980/golc/util"
)

// Compile time check to ensure SagemakerEndpoint satisfies the llm interface.
var _ golc.LLM = (*SagemakerEndpoint)(nil)

type Transformer interface {
	// Transforms the input to a format that model can accept
	// as the request Body. Should return bytes or seekable file
	// like object in the format specified in the content_type
	// request header.
	TransformInput(prompt string) ([]byte, error)

	// Transforms the output from the model to string that
	// the LLM class expects.
	TransformOutput(output []byte) (string, error)
}

type LLMContentHandler struct {
	// The MIME type of the input data passed to endpoint.
	contentType string

	// The MIME type of the response data returned from endpoint
	accept string

	transformer Transformer
}

func NewLLMContentHandler(contentType, accept string, transformer Transformer) *LLMContentHandler {
	return &LLMContentHandler{
		contentType: contentType,
		accept:      accept,
		transformer: transformer,
	}
}

func (ch *LLMContentHandler) ContentType() string {
	return ch.contentType
}

func (ch *LLMContentHandler) Accept() string {
	return ch.accept
}

func (ch *LLMContentHandler) TransformInput(prompt string) ([]byte, error) {
	return ch.transformer.TransformInput(prompt)
}

func (ch *LLMContentHandler) TransformOutput(output []byte) (string, error) {
	return ch.transformer.TransformOutput(output)
}

type SagemakerEndpoint struct {
	golc.Tokenizer
	client        *sagemakerruntime.Client
	endpointName  string
	contenHandler *LLMContentHandler
}

func NewSagemakerEndpoint(client *sagemakerruntime.Client, endpointName string, contenHandler *LLMContentHandler) (*SagemakerEndpoint, error) {
	se := &SagemakerEndpoint{
		Tokenizer:     tokenizer.NewSimple(),
		client:        client,
		endpointName:  endpointName,
		contenHandler: contenHandler,
	}

	return se, nil
}

func (se *SagemakerEndpoint) Generate(ctx context.Context, prompts []string) (*golc.LLMResult, error) {
	generations := [][]golc.Generation{}

	for _, prompt := range prompts {
		body, err := se.contenHandler.TransformInput(prompt)
		if err != nil {
			return nil, err
		}

		out, err := se.client.InvokeEndpoint(ctx, &sagemakerruntime.InvokeEndpointInput{
			EndpointName: aws.String(se.endpointName),
			ContentType:  aws.String(se.contenHandler.ContentType()),
			Accept:       aws.String(se.contenHandler.Accept()),
			Body:         body,
		})
		if err != nil {
			return nil, err
		}

		text, err := se.contenHandler.TransformOutput(out.Body)
		if err != nil {
			return nil, err
		}

		generations = append(generations, []golc.Generation{{
			Text: text,
		}})
	}

	return &golc.LLMResult{
		Generations: generations,
		LLMOutput:   map[string]any{},
	}, nil
}

func (se *SagemakerEndpoint) GeneratePrompt(ctx context.Context, promptValues []golc.PromptValue) (*golc.LLMResult, error) {
	prompts := util.Map(promptValues, func(value golc.PromptValue, _ int) string {
		return value.String()
	})

	return se.Generate(ctx, prompts)
}

func (se *SagemakerEndpoint) Predict(ctx context.Context, text string) (string, error) {
	result, err := se.Generate(ctx, []string{text})
	if err != nil {
		return "", err
	}

	return result.Generations[0][0].Text, nil
}

func (se *SagemakerEndpoint) PredictMessages(ctx context.Context, messages []golc.ChatMessage) (golc.ChatMessage, error) {
	text, err := golc.StringifyChatMessages(messages)
	if err != nil {
		return nil, err
	}

	prediction, err := se.Predict(ctx, text)
	if err != nil {
		return nil, err
	}

	return golc.NewAIChatMessage(prediction), nil
}
