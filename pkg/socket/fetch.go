package socket

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

// fetchRemoteStream loads the sockets to be monitored from a remote RESTful API endpoint. It returns the response body implementing [io.ReadCloser] for reading from and closing the stream.
//
// The url parameter must be a complete URL to a remote http/s server, including:
//   - Scheme (http:// or https://)
//   - Host (domain or IP)
//   - Optional port
//   - Optional path
//   - Optional query parameters
//
// Example url: http://api.example.com:5569/stream?query=variable
func fetchRemoteStream(url string, apiHeaderName string, apiHeaderValue string) (io.ReadCloser, error) {
	// try URL
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json")

	if apiHeaderName != "" && apiHeaderValue != "" {
		req.Header.Set(apiHeaderName, apiHeaderValue)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching sockets from remote source --- got %d (%s)", resp.StatusCode, resp.Status)
	}

	return resp.Body, nil
}

// fetchFileStream opens a file and returns [io.ReadCloser] for reading from and closing the stream.
func fetchFileStream(input string) (io.ReadCloser, error) {
	jsonFile, err := os.Open(input)
	if err != nil {
		return nil, err
	}

	return jsonFile, nil
}

// getStreamFromPath tries to open a stream to the socket source from the given input. It checks whether the input is a file path or url and then returns [io.ReadCloser] to read from and close the stream.
func getStreamFromPath(input string, apiHeaderName string, apiHeaderValue string) (io.ReadCloser, error) {
	// Check if input is a url
	match, err := regexp.MatchString("^(http|https)://", input)
	if err != nil {
		return nil, err
	}

	// If not, fetch from a file
	if !match {
		reader, err := fetchFileStream(input)
		if err != nil {
			return nil, err
		}
		return reader, nil
	}

	// Otherwise, fetch from the url
	reader, err := fetchRemoteStream(input, apiHeaderName, apiHeaderValue)
	if err != nil {
		return nil, err
	}

	return reader, nil
}
