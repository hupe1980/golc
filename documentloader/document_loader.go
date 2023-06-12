package documentloader

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/textsplitter"
)

type DocumentLoader interface {
	Load(context.Context) ([]golc.Document, error)
	LoadAndSplit(ctx context.Context, splitter textsplitter.TextSplitter)
}
