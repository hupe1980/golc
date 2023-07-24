package outputparser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure FencedCodeBlock satisfies the OutputParser interface.
var _ schema.OutputParser[any] = (*FencedCodeBlock)(nil)

// FencedCodeBlock represents a parser for extracting fenced code blocks from the output text.
type FencedCodeBlock struct {
	// The fence used to enclose the code block (e.g., "```" for code blocks in Markdown).
	fence string
}

// NewFencedCodeBlock creates a new instance of FencedCodeBlock with the specified fence.
func NewFencedCodeBlock(fence string) *FencedCodeBlock {
	return &FencedCodeBlock{
		fence: fence,
	}
}

// ParseResult parses the result of generation and returns the extracted fenced code blocks as a slice of strings.
func (p *FencedCodeBlock) ParseResult(result schema.Generation) (any, error) {
	return p.Parse(result.Text)
}

// Parse extracts fenced code blocks from the input text and returns them as a slice of strings.
// The fence used to enclose the code blocks is specified when creating the FencedCodeBlock instance.
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

// ParseWithPrompt is not used for this parser, so it simply calls Parse.
func (p *FencedCodeBlock) ParseWithPrompt(text string, prompt schema.PromptValue) (any, error) {
	return p.Parse(text)
}

// GetFormatInstructions returns a formatted string describing the expected format of the output.
// It instructs the user to enclose their response in a fenced code block using the specified fence.
func (p *FencedCodeBlock) GetFormatInstructions() string {
	return fmt.Sprintf("Your response should be enclosed in a fenced code block using three backticks (%s), e.g.: %s ls ```", p.fence, p.fence)
}

// Type returns the type identifier of the parser, which is "fenced_code_block".
func (p *FencedCodeBlock) Type() string {
	return "fenced_code_block"
}
