// Package retriever provides functionality for retrieving relevant documents using various services.
package retriever

import "net/http"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
