package documentloader

import (
	"context"
	"fmt"
	"os"

	"github.com/hupe1980/golc/integration"
	"github.com/hupe1980/golc/schema"
)

type Unstructured struct {
	client *integration.Unstructured
	file   *os.File
}

func NewUnstructured(apiKey string, file *os.File) *Unstructured {
	client := integration.NewUnstructured(apiKey)

	return &Unstructured{
		client: client,
		file:   file,
	}
}

func (l *Unstructured) Load(ctx context.Context) ([]schema.Document, error) {
	output, err := l.client.Partition(ctx, &integration.PartitionInput{
		File: l.file,
	})
	if err != nil {
		return nil, err
	}

	docs := []schema.Document{}

	textMap := make(map[int]string)
	metaMap := make(map[int]map[string]any)
	pages := []int{}

	for _, item := range output {
		page := item.Metadata.PageNumber

		if v, ok := textMap[page]; ok {
			textMap[page] = fmt.Sprintf("%s\n\n%s", v, item.Text)
			metaMap[page] = map[string]any{
				"page":      item.Metadata.PageNumber,
				"languages": item.Metadata.Languages,
				"filetype":  item.Metadata.Filetype,
				"filename":  item.Metadata.Filename,
			}
		} else {
			pages = append(pages, page)

			textMap[page] = item.Text
			metaMap[page] = map[string]any{
				"page":      item.Metadata.PageNumber,
				"languages": item.Metadata.Languages,
				"filetype":  item.Metadata.Filetype,
				"filename":  item.Metadata.Filename,
			}
		}
	}

	for _, p := range pages {
		docs = append(docs, schema.Document{
			PageContent: textMap[p],
			Metadata:    metaMap[p],
		})
	}

	return docs, nil
}

func (l *Unstructured) LoadAndSplit(ctx context.Context, splitter schema.TextSplitter) ([]schema.Document, error) {
	docs, err := l.Load(ctx)
	if err != nil {
		return nil, err
	}

	return splitter.SplitDocuments(docs)
}
