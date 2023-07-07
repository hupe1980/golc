package llm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemakerruntime"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
)

// Compile time check to ensure SagemakerEndpoint satisfies the LLM interface.
var _ schema.LLM = (*SagemakerEndpoint)(nil)

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

type SagemakerEndpointOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`
}

type SagemakerEndpoint struct {
	schema.Tokenizer
	client        *sagemakerruntime.Client
	endpointName  string
	contenHandler *LLMContentHandler
	opts          SagemakerEndpointOptions
}

func NewSagemakerEndpoint(client *sagemakerruntime.Client, endpointName string, contenHandler *LLMContentHandler, optFns ...func(o *SagemakerEndpointOptions)) (*SagemakerEndpoint, error) {
	opts := SagemakerEndpointOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Tokenizer == nil {
		var tErr error

		opts.Tokenizer, tErr = tokenizer.NewGPT2()
		if tErr != nil {
			return nil, tErr
		}
	}

	return &SagemakerEndpoint{
		Tokenizer:     opts.Tokenizer,
		client:        client,
		endpointName:  endpointName,
		contenHandler: contenHandler,
		opts:          opts,
	}, nil
}

// Generate generates text based on the provided prompt and options.
func (l *SagemakerEndpoint) Generate(ctx context.Context, prompt string, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	body, err := l.contenHandler.TransformInput(prompt)
	if err != nil {
		return nil, err
	}

	out, err := l.client.InvokeEndpoint(ctx, &sagemakerruntime.InvokeEndpointInput{
		EndpointName: aws.String(l.endpointName),
		ContentType:  aws.String(l.contenHandler.ContentType()),
		Accept:       aws.String(l.contenHandler.Accept()),
		Body:         body,
	})
	if err != nil {
		return nil, err
	}

	text, err := l.contenHandler.TransformOutput(out.Body)
	if err != nil {
		return nil, err
	}

	return &schema.ModelResult{
		Generations: []schema.Generation{{
			Text: text,
		}},
		LLMOutput: map[string]any{},
	}, nil
}

// Type returns the type of the model.
func (l *SagemakerEndpoint) Type() string {
	return "llm.SagemakerEndpoint"
}

// Verbose returns the verbosity setting of the model.
func (l *SagemakerEndpoint) Verbose() bool {
	return l.opts.CallbackOptions.Verbose
}

// Callbacks returns the registered callbacks of the model.
func (l *SagemakerEndpoint) Callbacks() []schema.Callback {
	return l.opts.CallbackOptions.Callbacks
}

// InvocationParams returns the parameters used in the model invocation.
func (l *SagemakerEndpoint) InvocationParams() map[string]any {
	return nil
}
