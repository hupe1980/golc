package documenttransformer

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/schema"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/sync/errgroup"
)

// Compile time check to ensure MetaDataTagger satisfies the DocumentTransformer interface.
var _ schema.DocumentTransformer = (*MetaDataTagger)(nil)

// MetaDataTaggerOptions represents the options for the MetaDataTagger.
type MetaDataTaggerOptions struct {
	MaxConcurrency int
}

// MetaDataTagger is a document transformer that adds metadata to documents using a tagging chain.
type MetaDataTagger struct {
	taggingChain *chain.Tagging
	opts         MetaDataTaggerOptions
}

// NewMetaDataTagger creates a new MetaDataTagger instance.
func NewMetaDataTagger(chatModel schema.ChatModel, tags any, optFns ...func(o *MetaDataTaggerOptions)) (*MetaDataTagger, error) {
	opts := MetaDataTaggerOptions{
		MaxConcurrency: 5,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	taggignChain, err := chain.NewTagging(chatModel, tags)
	if err != nil {
		return nil, err
	}

	return &MetaDataTagger{
		taggingChain: taggignChain,
		opts:         opts,
	}, nil
}

// Transform transforms a slice of documents by adding metadata using the tagging chain.
func (t *MetaDataTagger) Transform(ctx context.Context, docs []schema.Document) ([]schema.Document, error) {
	errs, errctx := errgroup.WithContext(ctx)

	errs.SetLimit(t.opts.MaxConcurrency)

	enrichedDocs := make([]schema.Document, len(docs))

	for i, d := range docs {
		i, d := i, d

		errs.Go(func() error {
			enrichedDocs[i] = d

			if enrichedDocs[i].PageContent == "" {
				return nil
			}

			result, err := golc.Call(errctx, t.taggingChain, schema.ChainValues{
				"input": d.PageContent,
			})
			if err != nil {
				return err
			}

			mr := map[string]any{}

			decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				TagName: "json",
				Result:  &mr,
			})
			if err != nil {
				return err
			}

			if err := decoder.Decode(result["output"]); err != nil {
				return err
			}

			if len(enrichedDocs[i].Metadata) > 0 {
				for k, v := range mr {
					enrichedDocs[i].Metadata[k] = v
				}
			} else {
				enrichedDocs[i].Metadata = mr
			}

			return nil
		})
	}

	if err := errs.Wait(); err != nil {
		return nil, err
	}

	return enrichedDocs, nil
}
