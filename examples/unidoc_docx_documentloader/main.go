package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc/documentloader"
	"github.com/hupe1980/golc/integration/unidoc"
)

func main() {
	parser, err := unidoc.New(os.Getenv("UNIDOC_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open("examples/docx_documentloader/document.docx")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	docx := documentloader.NewUniDocDOCX(parser, f)

	docs, err := docx.Load(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(docs[0].PageContent)
	fmt.Println(docs[0].Metadata)
}
