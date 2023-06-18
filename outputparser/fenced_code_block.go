package outputparser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure FencedCodeBlocks satisfies the OutputParser interface.
var _ schema.OutputParser[any] = (*FencedCodeBlock)(nil)

type FencedCodeBlock struct {
	fence string
}

func NewFencedCodeBlock(fence string) *FencedCodeBlock {
	return &FencedCodeBlock{
		fence: fence,
	}
}

func (p *FencedCodeBlock) Parse(text string) (any, error) {
	if !strings.Contains(text, p.fence) {
		return nil, fmt.Errorf("cannot parse output: %s", text)
	}

	codeBlocks := []string{}
	str := fmt.Sprintf("(?s)%s(.*?)\n\\s*```", p.fence)
	pattern := regexp.MustCompile(str)
	matches := pattern.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		matched := strings.TrimSpace(match[1])
		if matched != "" {
			lines := strings.Split(matched, "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					codeBlocks = append(codeBlocks, line)
				}
			}
		}
	}

	return codeBlocks, nil
}

func (p *FencedCodeBlock) ParseWithPrompt(text string, prompt schema.PromptValue) (any, error) {
	return p.Parse(text)
}

func (p *FencedCodeBlock) GetFormatInstructions() (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (p *FencedCodeBlock) Type() string {
	return "fenced-code-block-output-parser"
}
