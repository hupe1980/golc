// Package integration provides utilities for integration with external systems, services, and frameworks.
package integration

import "net/http"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
