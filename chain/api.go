package chain

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
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

// Compile time check to ensure API satisfies the Chain interface.
var _ schema.Chain = (*API)(nil)

// HTTPClient is an interface for making HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// VerifyURL is the function signature for verifying the API URL.
type VerifyURL func(url string) bool

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

// API represents a chain that makes API calls based on given API documentation and user questions.
//
// WARNING: The API chain has the potential to be susceptible to Server-Side Request Forgery (SSRF) attacks
// if not used carefully and securely. SSRF allows an attacker to manipulate the server into making unintended
// and unauthorized requests to internal or external resources, which can lead to potential security breaches
// and unauthorized access to sensitive information.
//
// To mitigate the risks associated with SSRF attacks, it is strongly advised to use the VerifyURL hook diligently.
// The VerifyURL hook should be implemented to validate and ensure that the generated URLs are restricted to authorized
// and safe resources only. By doing so, unauthorized access to sensitive resources can be prevented, and the application's
// security can be significantly enhanced.
//
// It is the responsibility of developers and administrators to ensure the secure usage of the API chain. We strongly recommend
// thorough testing, security reviews, and adherence to secure coding practices to protect against potential security threats,
// including SSRF and other vulnerabilities.
type API struct {
	apiRequestChain *LLM
	apiAnswerChain  *LLM
	apiDoc          string
	opts            APIOptions
}

// NewAPI creates a new instance of API with the given model, API documentation, and optional functions to set options.
func NewAPI(model schema.Model, apiDoc string, optFns ...func(o *APIOptions)) (*API, error) {
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

	apiRequestChain, err := NewLLM(model, prompt.NewTemplate(defaultAPIURLTemplate))
	if err != nil {
		return nil, err
	}

	apiAnswerChain, err := NewLLM(model, prompt.NewTemplate(defaultAPIAnswerTemplate))
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
	opts := schema.CallOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	question, err := inputs.GetString(c.opts.InputKey)
	if err != nil {
		return nil, err
	}

	if cbErr := opts.CallbackManger.OnText(ctx, &schema.TextManagerInput{
		Text: question,
	}); cbErr != nil {
		return nil, cbErr
	}

	apiURL, err := golc.SimpleCall(ctx, c.apiRequestChain, schema.ChainValues{
		"question": question,
		"apiDoc":   c.apiDoc,
	}, func(sco *golc.SimpleCallOptions) {
		sco.Callbacks = opts.CallbackManger.GetInheritableCallbacks()
		sco.ParentRunID = opts.CallbackManger.RunID()
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

	if cbErr := opts.CallbackManger.OnText(ctx, &schema.TextManagerInput{
		Text: apiURL,
	}); cbErr != nil {
		return nil, cbErr
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
	}, func(sco *golc.SimpleCallOptions) {
		sco.Callbacks = opts.CallbackManger.GetInheritableCallbacks()
		sco.ParentRunID = opts.CallbackManger.RunID()
	})
	if err != nil {
		return nil, err
	}

	answer = strings.TrimSpace(answer)

	if cbErr := opts.CallbackManger.OnText(ctx, &schema.TextManagerInput{
		Text: fmt.Sprintf("\nAnswer:\n%s", answer),
	}); cbErr != nil {
		return nil, cbErr
	}

	return schema.ChainValues{
		c.opts.OutputKey: answer,
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
