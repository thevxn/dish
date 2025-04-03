package socket

import (
	"fmt"
	"io"
	"net/http"
)

// fetchSocketsFromRemote loads the sockets to be monitored from a remote RESTful API endpoint. It returns the response body implementing [io.ReadCloser] for reading from and closing the stream.
//
// The url parameter must be a complete URL to a remote http/s server, including:
//   - Scheme (http:// or https://)
//   - Host (domain or IP)
//   - Optional port
//   - Optional path
//   - Optional query parameters
//
// Example url: http://api.example.com:5569/stream?query=variable
func fetchSocketsFromRemote(url string, cacheSockets bool, cacheDir string, cacheTTL uint, apiHeaderName string, apiHeaderValue string) (io.ReadCloser, error) {
	cacheFilePath := hashUrlToFilePath(url, cacheDir)
	respBody, err := attemptFetchFromNetwork(url, apiHeaderName, apiHeaderValue)
	if err != nil {
		return loadSocketsFromCache(cacheFilePath, cacheTTL)
	}

	// Caches response to the file
	if cacheSockets {
		err := saveSocketsToCache(cacheFilePath, cacheDir, respBody)
		if err != nil {
			return nil, err
		}

		return loadSocketsFromCache(cacheFilePath, cacheTTL)
	}

	// If we don't want to cache sockets to the file, return response body
	return respBody, nil
}

// attemptFetchFromNetwork tries to fetch sockets from the remote source.
func attemptFetchFromNetwork(url string, apiHeaderName string, apiHeaderValue string) (io.ReadCloser, error) {
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
