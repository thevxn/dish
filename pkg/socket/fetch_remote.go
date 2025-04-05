package socket

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

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
func fetchSocketsFromRemote(url string, cacheSockets bool, cacheDir string, cacheTTL uint, apiHeaderName string, apiHeaderValue string) (io.ReadCloser, error) {
	cacheFilePath := hashUrlToFilePath(url, cacheDir)

	if cacheSockets {
		cachedReader, cacheTime, err := loadCachedSockets(cacheFilePath, cacheTTL)
		// If cache is expired or fails to load, attempt to fetch fresh sockets
		if err != nil {
			if errors.Is(err, ErrExpiredCache) {
				log.Printf("cache expired for URL: %s. Attempting network fetch.", url)
			} else {
				log.Printf("failed to load cache for URL: %s. Attempting network fetch.", url)
			}

			// Fetch fresh sockets from network
			respBody, fetchErr := loadFreshSockets(url, apiHeaderName, apiHeaderValue)
			if fetchErr != nil {
				log.Printf("fetching socket list from remote API at %s failed: %v", url, fetchErr)

				// If the fetch fails and expired cache is not available, return the fetch error
				if err != ErrExpiredCache {
					return nil, fetchErr
				}

				// If the fetch fails and expired cache is available, return the expired cache and log a warning
				log.Printf("using expired cache from %s", cacheTime.Format(time.RFC3339))
				return cachedReader, nil
			}

			var buf bytes.Buffer
			_, err := buf.ReadFrom(respBody)
			if err != nil {
				return nil, fmt.Errorf("failed to copy response body: %w", err)
			}

			if err := saveSocketsToCache(cacheFilePath, cacheDir, buf.Bytes()); err != nil {
				log.Printf("failed to save fetched sockets to cache: %v", err)
			}

			return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
		}

		// Cache is valid (not expired, no error from file read)
		log.Println("loading sockets from cache...")
		return cachedReader, err
	}

	// If we do not want to cache sockets to the file, fetch from network
	return loadFreshSockets(url, apiHeaderName, apiHeaderValue)
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
