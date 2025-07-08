package alert

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const internalServerErrorResponse = "internal server error"

// SuccessStatusHTTPClient is a mock HTTP client implementation which returns HTTP Success (200) status responses.
type SuccessStatusHTTPClient struct{}

func (c *SuccessStatusHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("mocked Do response")),
	}, nil
}

func (c *SuccessStatusHTTPClient) Get(url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("mocked Get response")),
	}, nil
}

func (c *SuccessStatusHTTPClient) Post(
	url string,
	contentType string,
	body io.Reader,
) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("mocked Post response")),
	}, nil
}

// ErrorStatusHTTPClient is a mock HTTP client implementation which returns HTTP Internal Server Error (500) status responses.
type ErrorStatusHTTPClient struct{}

func (e *ErrorStatusHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(strings.NewReader(internalServerErrorResponse)),
		Request: &http.Request{
			URL: &url.URL{
				Host: "vxn.dev",
				Path: "/",
			},
		},
	}, nil
}

func (e *ErrorStatusHTTPClient) Get(url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(strings.NewReader(internalServerErrorResponse)),
	}, nil
}

func (e *ErrorStatusHTTPClient) Post(
	url, contentType string,
	body io.Reader,
) (*http.Response, error) {
	return &http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(strings.NewReader(internalServerErrorResponse)),
	}, nil
}

// FailureHTTPClient is a mock HTTP client implementation which simulates a failure to process the given request, returning nil as the response and an error.
type FailureHTTPClient struct{}

func (f *FailureHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("mocked Do error")
}

func (f *FailureHTTPClient) Get(url string) (*http.Response, error) {
	return nil, fmt.Errorf("mocked Get error")
}

func (f *FailureHTTPClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return nil, fmt.Errorf("mocked Post error")
}

// InvalidBodyReadCloser implements the [io.ReadCloser] interface and simulates an error when calling Read().
type InvalidBodyReadCloser struct{}

func (i *InvalidBodyReadCloser) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("invalid body")
}

func (i *InvalidBodyReadCloser) Close() error {
	return nil
}

// InvalidResponseBodyHTTPClient is a mock HTTP client implementation which simulates an invalid response body to trigger an error when trying to read it.
type InvalidResponseBodyHTTPClient struct{}

func (i *InvalidResponseBodyHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 500,
		Body:       &InvalidBodyReadCloser{},
	}, nil
}

func (i *InvalidResponseBodyHTTPClient) Get(url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: 500,
		Body:       &InvalidBodyReadCloser{},
	}, nil
}

func (i *InvalidResponseBodyHTTPClient) Post(
	url, contentType string,
	body io.Reader,
) (*http.Response, error) {
	return &http.Response{
		StatusCode: 500,
		Body:       &InvalidBodyReadCloser{},
	}, nil
}

// MockLogger is a mock implementation of the Logger interface with empty method implementations.
type MockLogger struct{}

func (l *MockLogger) Trace(v ...any)                 {}
func (l *MockLogger) Tracef(format string, v ...any) {}
func (l *MockLogger) Debug(v ...any)                 {}
func (l *MockLogger) Debugf(format string, v ...any) {}
func (l *MockLogger) Info(v ...any)                  {}
func (l *MockLogger) Infof(format string, v ...any)  {}
func (l *MockLogger) Warn(v ...any)                  {}
func (l *MockLogger) Warnf(format string, v ...any)  {}
func (l *MockLogger) Error(v ...any)                 {}
func (l *MockLogger) Errorf(format string, v ...any) {}
func (l *MockLogger) Panic(v ...any)                 {}
func (l *MockLogger) Panicf(format string, v ...any) {}
