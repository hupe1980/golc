package chain

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

const defaultAPIURLTemplate = `You are given the below API Documentation:
{{.apiDoc}}
Using this documentation, generate the full API url to call for answering the user question.
You should build the API url in order to get a response that is as short as possible, while still getting the necessary information to answer the question. Pay attention to deliberately exclude any unnecessary pieces of data in the API call.

Question:{{.question}}
API url:`

const defaultAPIAnswerTemplate = defaultAPIURLTemplate + `{{.apiURL}}

Here is the response from the API:

{{.apiResponse}}

Summarize this response to answer the original question.

Summary:`

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// VerifyURL is the function signature for verifying the API URL.
type VerifyURL func(url string) bool

// Compile time check to ensure API satisfies the Chain interface.
var _ schema.Chain = (*API)(nil)

type APIOptions struct {
	// CallbackOptions contains options for the chain callbacks.
	*schema.CallbackOptions

	// InputKey is the key to access the input value containing the user question.
	InputKey string

	// OutputKey is the key to access the output value containing the API response summary.
	OutputKey string

	// HTTPClient is the HTTP client used for making API requests.
	HTTPClient HTTPClient

	// Header is a map containing additional headers to be included in the API request.
	Header map[string]string

	// VerifyURL is a function used to verify the validity of the generated API URL before making the request.
	// It returns true if the URL is valid, false otherwise.
	VerifyURL VerifyURL
}

type API struct {
	apiRequestChain *LLM
	apiAnswerChain  *LLM
	apiDoc          string
	opts            APIOptions
}

func NewAPI(llm schema.Model, apiDoc string, optFns ...func(o *APIOptions)) (*API, error) {
	opts := APIOptions{
		InputKey:   "question",
		OutputKey:  "output",
		HTTPClient: http.DefaultClient,
		VerifyURL:  func(url string) bool { return true },
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	apiRequestChain, err := NewLLM(llm, prompt.NewTemplate(defaultAPIURLTemplate))
	if err != nil {
		return nil, err
	}

	apiAnswerChain, err := NewLLM(llm, prompt.NewTemplate(defaultAPIAnswerTemplate))
	if err != nil {
		return nil, err
	}

	return &API{
		apiRequestChain: apiRequestChain,
		apiAnswerChain:  apiAnswerChain,
		apiDoc:          apiDoc,
		opts:            opts,
	}, nil
}

// Call executes the api chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *API) Call(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
	question, ok := inputs[c.opts.InputKey].(string)
	if !ok {
		return nil, fmt.Errorf("%w: no value for inputKey %s", ErrInvalidInputValues, c.opts.InputKey)
	}

	apiURL, err := golc.SimpleCall(ctx, c.apiRequestChain, schema.ChainValues{
		"question": question,
		"apiDoc":   c.apiDoc,
	})
	if err != nil {
		return nil, err
	}

	apiURL = strings.TrimSpace(apiURL)
	if !strings.HasPrefix(apiURL, "https://") {
		apiURL = fmt.Sprintf("https://%s", apiURL)
	}

	if ok := c.opts.VerifyURL(apiURL); !ok {
		return nil, fmt.Errorf("invalid API URL: %s", apiURL)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	if c.opts.Header != nil {
		for k, v := range c.opts.Header {
			httpReq.Header.Set(k, v)
		}
	}

	res, err := c.opts.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	apiResponse, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	answer, err := golc.SimpleCall(ctx, c.apiAnswerChain, schema.ChainValues{
		"question":    question,
		"apiDoc":      c.apiDoc,
		"apiURL":      apiURL,
		"apiResponse": string(apiResponse),
	})
	if err != nil {
		return nil, err
	}

	return schema.ChainValues{
		c.opts.OutputKey: strings.TrimSpace(answer),
	}, nil
}

// Memory returns the memory associated with the chain.
func (c *API) Memory() schema.Memory {
	return nil
}

// Type returns the type of the chain.
func (c *API) Type() string {
	return "API"
}

// Verbose returns the verbosity setting of the chain.
func (c *API) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (c *API) Callbacks() []schema.Callback {
	return c.opts.CallbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *API) InputKeys() []string {
	return []string{c.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (c *API) OutputKeys() []string {
	return []string{c.opts.OutputKey}
}
