package alert

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// submitOptions holds optional parameters for submitting HTTP requests using handleSubmit.
type submitOptions struct {
	contentType string
	headers     map[string]string
}

// withContentType sets the provided contentType as the value of the Content-Type header.
func withContentType(contentType string) func(*submitOptions) {
	return func(opts *submitOptions) {
		opts.contentType = contentType
	}
}

// withHeader adds the provided key:value header pair to the request's HTTP headers.
func withHeader(key string, value string) func(*submitOptions) {
	return func(opts *submitOptions) {
		if opts.headers == nil {
			opts.headers = make(map[string]string)
		}
		opts.headers[key] = value
	}
}

// handleSubmit submits an HTTP request using the provided client and method to the specified url with the provided body (can be nil if no body is required).
//
// By default, the application/json Content-Type header is used. A different content type can be specified using the withContentType functional option.
// Custom header key:value pairs can be specified using the withHeader functional option.
//
// The response status code is checked and if it is not within the range of success codes (2xx), the response body is logged and an error with the received status code is returned.
func handleSubmit(client HTTPClient, method string, url string, body io.Reader, opts ...func(*submitOptions)) error {
	// Default options
	options := submitOptions{
		contentType: "application/json",
		headers:     make(map[string]string),
	}

	// Apply provided options to the defaults
	for _, opt := range opts {
		opt(&options)
	}

	// Prepare the request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	// Set content type
	req.Header.Set("Content-Type", options.contentType)

	// Apply provided custom headers
	for k, v := range options.headers {
		req.Header.Set(k, v)
	}

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// If status code is not within <200, 299>, log the body and return an error with the received status code
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("error reading response body: %v", err)
		} else {
			log.Printf("response from %s: %s", url, string(body))
		}

		return fmt.Errorf("unexpected response code received (expected: 200-299, got: %d)", res.StatusCode)
	}

	return nil
}
