package socket

import (
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"testing"

	"go.vxn.dev/dish/pkg/config"
	"go.vxn.dev/dish/pkg/testhelpers"
)

func TestFetchSocketsFromRemote(t *testing.T) {
	apiHeaderName := "Authorization"
	apiHeaderValue := "Bearer xyzzzzzzz"
	mockServer := testhelpers.NewMockServer(t, apiHeaderName, apiHeaderValue, testhelpers.TestSocketList, http.StatusOK)

	newConfig := func(source string, useCache bool, ttl uint) *config.Config {
		// Temp cache directory needs to be created and specified for each test separately
		// See the range tests below
		return &config.Config{
			Source:             source,
			ApiCacheSockets:    useCache,
			ApiCacheTTLMinutes: ttl,
			ApiHeaderName:      apiHeaderName,
			ApiHeaderValue:     apiHeaderValue,
		}
	}

	tests := []struct {
		name          string
		cfg           *config.Config
		expectedError bool
	}{
		{"Fetch With Valid Cache", newConfig(mockServer.URL, true, 10), false},
		{"Fetch With Expired Cache", newConfig(mockServer.URL, true, 0), false},
		{"Fetch Without Caching", newConfig(mockServer.URL, false, 0), false},
		{"Invalid URL Without Cache", newConfig("http://badurl.com", false, 0), true},
		{"Invalid URL With Cache", newConfig("http://badurl.com", true, 0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Specify temp cache file & directory for each test separately
			// This fixes open file handles preventing the tests from succeeding on Windows
			filePath := testhelpers.TestFile(t, "randomhash.json", []byte(testhelpers.TestSocketList))
			cacheDir := filepath.Dir(filePath)
			tt.cfg.ApiCacheDirectory = cacheDir

			resp, err := fetchSocketsFromRemote(tt.cfg)
			if tt.expectedError {
				if err == nil || errors.Is(err, ErrExpiredCache) {
					t.Errorf("expected error, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			readBytes, err := io.ReadAll(resp)
			if err != nil {
				t.Fatalf("failed to read from response: %v", err)
			}

			if string(readBytes) != testhelpers.TestSocketList {
				t.Errorf("expected %s, got %s", testhelpers.TestSocketList, string(readBytes))
			}
		})
	}
}
