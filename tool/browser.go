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

func (t *CurrentPage) Name() string {
	return "CurrentPage"
}

func (t *CurrentPage) Description() string {
	return `Returns the URL of the current page.`
}

func (t *CurrentPage) ArgsType() reflect.Type {
	return reflect.TypeOf("") // string
}

func (t *CurrentPage) Run(ctx context.Context, input any) (string, error) {
	page, err := getCurrentPage(t.browser)
	if err != nil {
		return "", err
	}

	return page.URL(), nil
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

func (t *ExtractText) Name() string {
	return "ExtractText"
}

func (t *ExtractText) Description() string {
	return `Extract all the text on the current webpage.`
}

func (t *ExtractText) ArgsType() reflect.Type {
	return reflect.TypeOf("") // string
}

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

func (t *NavigateBrowser) Name() string {
	return "NavigateBrowser"
}

func (t *NavigateBrowser) Description() string {
	return `Navigate a browser to the specified URL.`
}

func (t *NavigateBrowser) ArgsType() reflect.Type {
	return reflect.TypeOf("") // string
}

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
