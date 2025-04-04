package alert

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// successStatusHTTPClient is a mock HTTP client implementation which returns HTTP Success (200) status responses.
type successStatusHTTPClient struct{}

func (c *successStatusHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("mocked Do response")),
	}, nil
}

func (c *successStatusHTTPClient) Get(url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("mocked Get response")),
	}, nil
}

func (c *successStatusHTTPClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("mocked Post response")),
	}, nil
}

// errorStatusHTTPClient is a mock HTTP client implementation which returns HTTP Internal Server Error (500) status responses.
type errorStatusHTTPClient struct{}

func (e *errorStatusHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(strings.NewReader("internal server error")),
	}, nil
}

func (e *errorStatusHTTPClient) Get(url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(strings.NewReader("internal server error")),
	}, nil
}

func (e *errorStatusHTTPClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return &http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(strings.NewReader("internal server error")),
	}, nil
}

// failureHTTPClient is a mock HTTP client implementation which simulates a failure to process the given request, returning nil as the response and an error.
type failureHTTPClient struct{}

func (f *failureHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("mocked Do error")
}

func (f *failureHTTPClient) Get(url string) (*http.Response, error) {
	return nil, fmt.Errorf("mocked Get error")
}

func (f *failureHTTPClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return nil, fmt.Errorf("mocked Post error")
}

// invalidBodyReadCloser implements the [io.ReadCloser] interface and simulates an error when calling Read().
type invalidBodyReadCloser struct{}

func (i *invalidBodyReadCloser) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("invalid body")
}

func (i *invalidBodyReadCloser) Close() error {
	return nil
}

// invalidResponseBodyHTTPClient is a mock HTTP client implementation which simulates an invalid response body to trigger an error when trying to read it.
type invalidResponseBodyHTTPClient struct{}

func (i *invalidResponseBodyHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       &invalidBodyReadCloser{},
	}, nil
}

func (i *invalidResponseBodyHTTPClient) Get(url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       &invalidBodyReadCloser{},
	}, nil
}

func (i *invalidResponseBodyHTTPClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       &invalidBodyReadCloser{},
	}, nil
}
