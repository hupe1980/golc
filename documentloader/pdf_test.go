package documentloader

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPDF(t *testing.T) {
	t.Run("Load PDF", func(t *testing.T) {
		// Load the test PDF file from the testdata folder
		file, err := os.Open("testdata/testfile.pdf")
		require.NoError(t, err)

		defer file.Close()

		finfo, err := file.Stat()
		require.NoError(t, err)

		pdfLoader, err := NewPDF(file, finfo.Size())
		require.NoError(t, err)

		docs, err := pdfLoader.Load(context.Background())
		require.NoError(t, err)
		require.Equal(t, 3, len(docs))
		require.Equal(t, "Page 1: Text text text", docs[0].PageContent)
		require.Equal(t, "Page 2: Text text text", docs[1].PageContent)
		require.Equal(t, "Page 3: Text text text", docs[2].PageContent)
	})

	t.Run("Load PDF with Password", func(t *testing.T) {
		// Load the test PDF file from the testdata folder
		file, err := os.Open("testdata/testfile_password.pdf")
		require.NoError(t, err)

		defer file.Close()

		finfo, err := file.Stat()
		require.NoError(t, err)

		pdfLoader, err := NewPDF(file, finfo.Size(), func(o *PDFOptions) {
			o.Password = "Secret"
		})
		require.NoError(t, err)

		docs, err := pdfLoader.Load(context.Background())
		require.NoError(t, err)
		require.Equal(t, 3, len(docs))
		require.Equal(t, "Page 1: Text text text", docs[0].PageContent)
		require.Equal(t, "Page 2: Text text text", docs[1].PageContent)
		require.Equal(t, "Page 3: Text text text", docs[2].PageContent)
	})

	t.Run("Load PDF with MaxPages", func(t *testing.T) {
		// Load the test PDF file from the testdata folder
		file, err := os.Open("testdata/testfile.pdf")
		require.NoError(t, err)

		defer file.Close()

		finfo, err := file.Stat()
		require.NoError(t, err)

		pdfLoader, err := NewPDF(file, finfo.Size(), func(o *PDFOptions) {
			o.MaxPages = 2
		})
		require.NoError(t, err)

		docs, err := pdfLoader.Load(context.Background())
		require.NoError(t, err)
		require.Equal(t, 2, len(docs))
		require.Equal(t, "Page 1: Text text text", docs[0].PageContent)
		require.Equal(t, "Page 2: Text text text", docs[1].PageContent)
	})

	t.Run("Load PDF with StartPage", func(t *testing.T) {
		// Load the test PDF file from the testdata folder
		file, err := os.Open("testdata/testfile.pdf")
		require.NoError(t, err)

		defer file.Close()

		finfo, err := file.Stat()
		require.NoError(t, err)

		pdfLoader, err := NewPDF(file, finfo.Size(), func(o *PDFOptions) {
			o.StartPage = 2
		})
		require.NoError(t, err)

		docs, err := pdfLoader.Load(context.Background())
		require.NoError(t, err)
		require.Equal(t, 2, len(docs))
		require.Equal(t, "Page 2: Text text text", docs[0].PageContent)
		require.Equal(t, "Page 3: Text text text", docs[1].PageContent)
	})

	t.Run("Load PDF with StartPage out of range", func(t *testing.T) {
		// Load the test PDF file from the testdata folder
		file, err := os.Open("testdata/testfile.pdf")
		require.NoError(t, err)

		defer file.Close()

		finfo, err := file.Stat()
		require.NoError(t, err)

		pdfLoader, err := NewPDF(file, finfo.Size(), func(o *PDFOptions) {
			o.StartPage = 4
		})
		require.NoError(t, err)

		_, err = pdfLoader.Load(context.Background())
		require.EqualError(t, err, "startpage out of page range: 1-3")
	})
}
