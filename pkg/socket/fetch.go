package socket

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"go.vxn.dev/dish/pkg/config"
	"go.vxn.dev/dish/pkg/logger"
)

// fetchHandler provides methods to fetch sockets either from a file or from a remote API source.
type fetchHandler struct {
	logger logger.Logger
}

// NewFetchHandler creates a new instance of fetchHandler.
func NewFetchHandler(l logger.Logger) *fetchHandler {
	return &fetchHandler{
		logger: l,
	}
}

// fetchSocketsFromFile opens a file and returns [io.ReadCloser] for reading from the stream.
func (f *fetchHandler) fetchSocketsFromFile(config *config.Config) (io.ReadCloser, error) {
	file, err := os.Open(config.Source)
	if err != nil {
		return nil, err
	}

	f.logger.Debugf("fetching sockets from file (%s)", config.Source)

	return file, nil
}

// copyBody copies the provided response body to the provided buffer. The body is closed.
func (f *fetchHandler) copyBody(body io.ReadCloser, buf *bytes.Buffer) error {
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
func (f *fetchHandler) fetchSocketsFromRemote(config *config.Config) (io.ReadCloser, error) {
	cacheFilePath := hashUrlToFilePath(config.Source, config.ApiCacheDirectory)

	// If we do not want to cache sockets to the file, fetch from network
	if !config.ApiCacheSockets {
		return f.loadFreshSockets(config)
	}

	// If cache is enabled, try to load sockets from it first
	cachedReader, cacheTime, err := loadCachedSockets(cacheFilePath, config.ApiCacheTTLMinutes)
	// If cache is expired or fails to load, attempt to fetch fresh sockets
	if err != nil {
		f.logger.Warnf(
			"cache unavailable for URL: %s (reason: %v); attempting network fetch",
			config.Source,
			err,
		)

		// Fetch fresh sockets from network
		respBody, fetchErr := f.loadFreshSockets(config)
		if fetchErr != nil {
			// If the fetch fails and expired cache is not available, return the fetch error
			if err != ErrExpiredCache {
				return nil, fetchErr
			}
			// If the fetch fails and expired cache is available, return the expired cache and log a warning
			f.logger.Errorf(
				"fetching socket list from remote API at %s failed: %v.",
				config.Source,
				fetchErr,
			)
			f.logger.Warnf("using expired cache from %s", cacheTime.Format(time.RFC3339))

			return cachedReader, nil
		} else {
			f.logger.Infof("socket list fetched from %s", config.Source)
		}

		var buf bytes.Buffer
		err = f.copyBody(respBody, &buf)
		if err != nil {
			return nil, fmt.Errorf("failed to copy response body: %w", err)
		}

		if err := saveSocketsToCache(cacheFilePath, config.ApiCacheDirectory, buf.Bytes()); err != nil {
			f.logger.Warnf("failed to save fetched sockets to cache: %v", err)
		}

		return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
	}

	// Cache is valid (not expired, no error from file read)
	f.logger.Info("socket list fetched from cache")
	return cachedReader, err
}

// loadFreshSockets fetches fresh sockets from the remote source.
func (f *fetchHandler) loadFreshSockets(config *config.Config) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, config.Source, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json")

	if config.ApiHeaderName != "" && config.ApiHeaderValue != "" {
		req.Header.Set(config.ApiHeaderName, config.ApiHeaderValue)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"failed to fetch sockets from remote source --- got %d (%s)",
			resp.StatusCode,
			resp.Status,
		)
	}

	return resp.Body, nil
}
