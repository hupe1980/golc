package toolkit

import (
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tool"
	"github.com/playwright-community/playwright-go"
)

// Browser represents a collection of schema.Tool objects that enable interaction with a browser.
type Browser struct {
	tools []schema.Tool
}

// NewBrowser creates a new Browser object from the given playwright.Browser instance.
// It initializes various schema.Tool objects that facilitate interactions with the browser.
func NewBrowser(browser playwright.Browser) (*Browser, error) {
	tools := []schema.Tool{
		tool.NewCurrentPage(browser),
		tool.NewNavigateBrowser(browser),
		tool.NewExtractText(browser),
	}

	return &Browser{
		tools: tools,
	}, nil
}

// Tools returns the list of schema.Tool objects associated with the Browser.
func (tk *Browser) Tools() []schema.Tool {
	return tk.tools
}
