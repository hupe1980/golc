package outputparser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hupe1980/golc"
)

// Compile time check to ensure BashOutputParse satisfies the output_parser interface.
var _ golc.OutputParser[any] = (*BashOutputParser)(nil)

type BashOutputParser struct{}

func NewBashOutputParser() *BashOutputParser {
	return &BashOutputParser{}
}

func (p *BashOutputParser) Parse(text string) (any, error) {
	if !strings.Contains(text, "```bash") {
		return nil, fmt.Errorf("cannot parse bash output: %s", text)
	}

	codeBlocks := []string{}
	pattern := regexp.MustCompile("(?s)```bash(.*?)\n\\s*```")
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

func (p *BashOutputParser) ParseWithPrompt(text string, prompt golc.PromptValue) (any, error) {
	return p.Parse(text)
}

func (p *BashOutputParser) GetFormatInstructions() (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (p *BashOutputParser) Type() string {
	return "bash-output-parser"
}
