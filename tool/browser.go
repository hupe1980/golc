package tool

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/util"
	"github.com/playwright-community/playwright-go"
)

// Compile time check to ensure CurrentPage satisfies the Tool interface.
var _ schema.Tool = (*CurrentPage)(nil)

type CurrentPage struct {
	browser playwright.Browser
}

func NewCurrentPage(browser playwright.Browser) *CurrentPage {
	return &CurrentPage{
		browser: browser,
	}
}

// Name returns the name of the tool.
func (t *CurrentPage) Name() string {
	return "CurrentPage"
}

// Description returns the description of the tool.
func (t *CurrentPage) Description() string {
	return `Returns the URL of the current page.`
}

// ArgsType returns the type of the input argument expected by the tool.
func (t *CurrentPage) ArgsType() reflect.Type {
	return reflect.TypeOf("") // string
}

// Run executes the tool with the given input and returns the output.
func (t *CurrentPage) Run(ctx context.Context, input any) (string, error) {
	page, err := getCurrentPage(t.browser)
	if err != nil {
		return "", err
	}

	return page.URL(), nil
}

// Verbose returns the verbosity setting of the tool.
func (t *CurrentPage) Verbose() bool {
	return false
}

// Callbacks returns the registered callbacks of the tool.
func (t *CurrentPage) Callbacks() []schema.Callback {
	return nil
}

// Compile time check to ensure ExtractText satisfies the Tool interface.
var _ schema.Tool = (*ExtractText)(nil)

type ExtractText struct {
	browser playwright.Browser
}

func NewExtractText(browser playwright.Browser) *ExtractText {
	return &ExtractText{
		browser: browser,
	}
}

// Name returns the name of the tool.
func (t *ExtractText) Name() string {
	return "ExtractText"
}

// Description returns the description of the tool.
func (t *ExtractText) Description() string {
	return `Extract all the text on the current webpage.`
}

// ArgsType returns the type of the input argument expected by the tool.
func (t *ExtractText) ArgsType() reflect.Type {
	return reflect.TypeOf("") // string
}

// Run executes the tool with the given input and returns the output.
func (t *ExtractText) Run(ctx context.Context, input any) (string, error) {
	page, err := getCurrentPage(t.browser)
	if err != nil {
		return "", err
	}

	html, err := page.Content()
	if err != nil {
		return "", err
	}

	return util.ParseHTMLAndGetStrippedStrings(html)
}

// Verbose returns the verbosity setting of the tool.
func (t *ExtractText) Verbose() bool {
	return false
}

// Callbacks returns the registered callbacks of the tool.
func (t *ExtractText) Callbacks() []schema.Callback {
	return nil
}

// Compile time check to ensure NavigateBrowser satisfies the Tool interface.
var _ schema.Tool = (*NavigateBrowser)(nil)

type NavigateBrowser struct {
	browser playwright.Browser
}

func NewNavigateBrowser(browser playwright.Browser) *NavigateBrowser {
	return &NavigateBrowser{
		browser: browser,
	}
}

// Name returns the name of the tool.
func (t *NavigateBrowser) Name() string {
	return "NavigateBrowser"
}

// Description returns the description of the tool.
func (t *NavigateBrowser) Description() string {
	return `Navigate a browser to the specified URL.`
}

// ArgsType returns the type of the input argument expected by the tool.
func (t *NavigateBrowser) ArgsType() reflect.Type {
	return reflect.TypeOf("") // string
}

// Run executes the tool with the given input and returns the output.
func (t *NavigateBrowser) Run(ctx context.Context, input any) (string, error) {
	url, ok := input.(string)
	if !ok {
		return "", errors.New("illegal input type")
	}

	if !strings.HasPrefix(url, "http") {
		url = fmt.Sprintf("https://%s", url)
	}

	page, err := getCurrentPage(t.browser)
	if err != nil {
		return "", err
	}

	res, err := page.Goto(url)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Navigating to %s returned status code %d", url, res.Status()), nil
}

// Verbose returns the verbosity setting of the tool.
func (t *NavigateBrowser) Verbose() bool {
	return false
}

// Callbacks returns the registered callbacks of the tool.
func (t *NavigateBrowser) Callbacks() []schema.Callback {
	return nil
}

func getCurrentPage(browser playwright.Browser) (playwright.Page, error) {
	if len(browser.Contexts()) == 0 {
		context, err := browser.NewContext()
		if err != nil {
			return nil, err
		}

		return context.NewPage()
	}

	context := browser.Contexts()[0]

	pages := context.Pages()
	if len(pages) == 0 {
		return context.NewPage()
	}

	return context.Pages()[len(pages)-1], nil
}
