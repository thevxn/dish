package socket

import (
	"bytes"
	"fmt"
	"io"
	"log"
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

	// Cache feature
	if cacheSockets {
		reader, err := loadSocketsFromCache(cacheFilePath, cacheTTL)
		// Attempt to fetch fresh sockets if cache is expired
		if err != nil {
			if err == ErrExpiredCache {
				log.Printf("Cache expired for URL: %s. Attempting network fetch.\n", url)
			} else {
				log.Printf("Failed to load cache for URL: %s. Attempting network fetch.\n", url)
			}

			// Fetch fresh sockets from network (returns expired cache on fail)
			respBody, fetchErr := attemptFetchFromNetwork(url, apiHeaderName, apiHeaderValue)
			if fetchErr != nil {
				if err != ErrExpiredCache {
					return nil, fetchErr
				}

				log.Printf("Network fetch failed for URL: %s. Using expired cache. Error: %v\n", url, err)
				return reader, nil
			}

			// TeeReader clones respBody reader which we pass to save func and return
			var buffer bytes.Buffer
			teeReader := io.TeeReader(respBody, &buffer)

			if err := saveSocketsToCache(cacheFilePath, cacheDir, io.NopCloser(teeReader)); err != nil {
				return nil, fmt.Errorf("failed to save fetched sockets to cache: %w", err)
			}
			log.Println("Successfully cached sockets from", url)

			return io.NopCloser(&buffer), nil
		}

		// Cache is valid (not expired, no error from file read)
		return reader, err
	}

	// If we do not want to cache sockets to the file, fetch from network
	return attemptFetchFromNetwork(url, apiHeaderName, apiHeaderValue)
}

// attemptFetchFromNetwork tries to fetch sockets from the remote source.
func attemptFetchFromNetwork(url string, apiHeaderName string, apiHeaderValue string) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json")

	if apiHeaderName != "" && apiHeaderValue != "" {
		req.Header.Set(apiHeaderName, apiHeaderValue)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch sockets from remote source --- got %d (%s)", resp.StatusCode, resp.Status)
	}

	return resp.Body, nil
}
