package unidoc

import (
	"io"

	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/document"
)

type Document interface {
	ExtractText() *document.DocText
}

type UniDoc struct{}

func New(apiKey string) (*UniDoc, error) {
	err := license.SetMeteredKey(apiKey)
	if err != nil {
		return nil, err
	}

	return &UniDoc{}, nil
}

func (u *UniDoc) ReadDocument(r io.ReaderAt, size int64) (Document, error) {
	return document.Read(r, size)
}
