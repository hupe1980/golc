package unidoc

import (
	"os"

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

func (u *UniDoc) ReadDocument(f *os.File) (Document, error) {
	finfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return document.Read(f, finfo.Size())
}
