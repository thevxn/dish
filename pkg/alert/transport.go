package alert

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.vxn.dev/dish/pkg/logger"
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

// handleSubmit submits an HTTP request using the provided client and method to the specified url with the provided body (can be nil if no body is required) and returns the response.
//
// By default, the application/json Content-Type header is used. A different content type can be specified using the withContentType functional option.
// Custom header key:value pairs can be specified using the withHeader functional option.
func handleSubmit(client HTTPClient, method string, url string, body io.Reader, opts ...func(*submitOptions)) (*http.Response, error) {
	// Default options
	options := submitOptions{
		contentType: "application/json",
		headers:     make(map[string]string),
	}

	// Apply provided options, if any, to the defaults
	for _, opt := range opts {
		opt(&options)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", options.contentType)

	// Apply provided custom headers, if any
	for k, v := range options.headers {
		req.Header.Set(k, v)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// handleRead reads an HTTP response, ensures the status code is within the expected <200, 299> range and if not, logs the response body.
func handleRead(res *http.Response, logger logger.Logger) error {
	defer func() {
		if err := res.Body.Close(); err != nil {
			logger.Errorf("failed to close response body: %v", err)
		}
	}()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			logger.Errorf("error reading response body: %v", err)
		} else {
			logger.Warnf("response from %s: %s", res.Request.URL, string(body))
		}

		return fmt.Errorf("unexpected response code received (expected: 200-299, got: %d)", res.StatusCode)
	}

	return nil
}
