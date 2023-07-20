// Package textsplitter provides utilities for splitting and processing text.
package textsplitter

import (
	"log"
	"regexp"
	"strings"

	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/util"
)

type SplitTextFunc func(text string) []string

type LengthFunc func(text string) int

type Options struct {
	ChunkSize     int
	ChunkOverlap  int
	KeepSeparator bool
	LengthFunc    LengthFunc
}

type BaseTextSplitter struct {
	splitTextFunc SplitTextFunc
	opts          Options
}

func NewBaseTextSplitter(splitTextFunc SplitTextFunc, optFns ...func(o *Options)) *BaseTextSplitter {
	opts := Options{
		ChunkSize:     4000,
		ChunkOverlap:  200,
		KeepSeparator: false,
		LengthFunc: func(text string) int {
			return len(text)
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &BaseTextSplitter{
		splitTextFunc: splitTextFunc,
		opts:          opts,
	}
}

func (ts *BaseTextSplitter) CreateDocuments(texts []string, metadatas []map[string]any) ([]schema.Document, error) {
	docs := []schema.Document{}

	for i, text := range texts {
		for _, chunk := range ts.splitTextFunc(text) {
			metadata := util.CopyMap(metadatas[i])
			docs = append(docs, schema.Document{
				PageContent: chunk,
				Metadata:    metadata,
			})
		}
	}

	return docs, nil
}

func (ts *BaseTextSplitter) SplitDocuments(docs []schema.Document) ([]schema.Document, error) {
	texts := []string{}
	metadatas := []map[string]any{}

	for _, doc := range docs {
		if doc.PageContent == "" {
			continue
		}

		texts = append(texts, doc.PageContent)
		metadatas = append(metadatas, doc.Metadata)
	}

	return ts.CreateDocuments(texts, metadatas)
}

func (ts *BaseTextSplitter) mergeSplits(splits []string, separator string) []string {
	separatorLen := ts.opts.LengthFunc(separator)
	docs := make([]string, 0)
	currentDoc := make([]string, 0)
	total := 0

	for _, d := range splits {
		lenD := ts.opts.LengthFunc(d)
		if total+lenD+(separatorLen*func() int {
			if len(currentDoc) > 0 {
				return 1
			}
			return 0
		}()) > ts.opts.ChunkSize {
			if total > ts.opts.ChunkSize {
				log.Println("Created a chunk of size", total, "which is longer than the specified", ts.opts.ChunkSize)
			}

			if len(currentDoc) > 0 {
				doc := ts.joinDocs(currentDoc, separator)
				if doc != nil {
					docs = append(docs, *doc)
				}

				for total > ts.opts.ChunkOverlap || (total+lenD+(separatorLen*func() int {
					if len(currentDoc) > 0 {
						return 1
					}
					return 0
				}()) > ts.opts.ChunkSize && total > 0) {
					total -= len(currentDoc[0]) + (separatorLen * func() int {
						if len(currentDoc) > 1 {
							return 1
						}
						return 0
					}())
					currentDoc = currentDoc[1:]
				}
			}
		}

		currentDoc = append(currentDoc, d)

		total += lenD + (separatorLen * func() int {
			if len(currentDoc) > 1 {
				return 1
			}
			return 0
		}())
	}

	doc := ts.joinDocs(currentDoc, separator)
	if doc != nil {
		docs = append(docs, *doc)
	}

	return docs
}

func (ts *BaseTextSplitter) joinDocs(docs []string, separator string) *string {
	text := strings.Join(docs, separator)
	text = strings.TrimSpace(text)

	if text == "" {
		return nil
	}

	return &text
}

// splitTextWithRegex splits the given text using the specified separator with optional
// inclusion of the separator itself.
func splitTextWithRegex(text string, separator string, keepSeparator bool) []string {
	var splits []string

	if separator != "" {
		if keepSeparator {
			splits = regexp.MustCompile("("+separator+")").Split(text, -1)
			processedSplits := make([]string, 0, len(splits))

			for i := 1; i < len(splits); i += 2 {
				processedSplits = append(processedSplits, splits[i]+splits[i+1])
			}

			if len(splits)%2 == 0 {
				processedSplits = append(processedSplits, splits[len(splits)-1])
			}

			splits = append([]string{splits[0]}, processedSplits...)
		} else {
			splits = strings.Split(text, separator)
		}
	} else {
		splits = strings.Split(text, "")
	}

	filteredSplits := make([]string, 0)

	for _, s := range splits {
		if s != "" {
			filteredSplits = append(filteredSplits, s)
		}
	}

	return filteredSplits
}
