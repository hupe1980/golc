package toolkit

import (
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tool"
	"github.com/playwright-community/playwright-go"
)

type Browser struct {
	tools []schema.Tool
}

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

func (tk *Browser) Tools() []schema.Tool {
	return tk.tools
}
