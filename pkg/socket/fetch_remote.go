package socket

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go.vxn.dev/dish/pkg/config"
)

// copyBody copies the provided response body to the provided buffer. The body is closed.
func copyBody(body io.ReadCloser, buf *bytes.Buffer) error {
	defer body.Close()

	_, err := buf.ReadFrom(body)
	return err
}

// fetchSocketsFromRemote loads the sockets to be monitored from a remote RESTful API endpoint. It returns the response body implementing [io.ReadCloser] for reading from and closing the stream.
//
// It uses a local cache if enabled and falls back to the network if the cache is not present or expired. If the network request fails and expired cache is available, it will be used.
//
// The url parameter must be a complete URL to a remote http/s server, including:
//   - Scheme (http:// or https://)
//   - Host (domain or IP)
//   - Optional port
//   - Optional path
//   - Optional query parameters
//
// Example url: http://api.example.com:5569/stream?query=variable
func fetchSocketsFromRemote(config *config.Config) (io.ReadCloser, error) {
	cacheFilePath := hashUrlToFilePath(config.Source, config.ApiCacheDirectory)

	// If we do not want to cache sockets to the file, fetch from network
	if !config.ApiCacheSockets {
		return loadFreshSockets(config.Source, config.ApiHeaderName, config.ApiHeaderValue)
	}

	// If cache is enabled, try to load sockets from it first
	cachedReader, cacheTime, err := loadCachedSockets(cacheFilePath, config.ApiCacheTTLMinutes)
	// If cache is expired or fails to load, attempt to fetch fresh sockets
	if err != nil {
		log.Printf("cache unavailable for URL: %s (reason: %v); attempting network fetch", config.Source, err)

		// Fetch fresh sockets from network
		respBody, fetchErr := loadFreshSockets(config.Source, config.ApiHeaderName, config.ApiHeaderValue)
		if fetchErr != nil {
			log.Printf("fetching socket list from remote API at %s failed: %v", config.Source, fetchErr)

			// If the fetch fails and expired cache is not available, return the fetch error
			if err != ErrExpiredCache {
				return nil, fetchErr
			}
			// If the fetch fails and expired cache is available, return the expired cache and log a warning
			log.Printf("using expired cache from %s", cacheTime.Format(time.RFC3339))
			return cachedReader, nil
		}

		var buf bytes.Buffer
		err = copyBody(respBody, &buf)
		if err != nil {
			return nil, fmt.Errorf("failed to copy response body: %w", err)
		}

		if err := saveSocketsToCache(cacheFilePath, config.ApiCacheDirectory, buf.Bytes()); err != nil {
			log.Printf("failed to save fetched sockets to cache: %v", err)
		}

		return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
	}

	// Cache is valid (not expired, no error from file read)
	log.Println("loading sockets from cache...")
	return cachedReader, err
}

// loadFreshSockets fetches fresh sockets from the remote source.
func loadFreshSockets(url string, apiHeaderName string, apiHeaderValue string) (io.ReadCloser, error) {
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
