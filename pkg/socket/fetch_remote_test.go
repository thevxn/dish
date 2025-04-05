package socket

import (
	"io"
	"net/http"
	"path/filepath"
	"testing"

	"go.vxn.dev/dish/pkg/testhelpers"
)

func TestFetchSocketsFromRemote(t *testing.T) {
	apiHeaderName := "Authorization"
	apiHeaderValue := "Bearer xyzzzzzzz"

	mockServer := testhelpers.NewMockServer(t, apiHeaderName, apiHeaderValue, testhelpers.TestSocketList, http.StatusOK)
	filePath := testhelpers.TestFile(t, "randomhash.json", []byte(testhelpers.TestSocketList))
	cacheDir := filepath.Dir(filePath)

	t.Run("Fetch with valid cache", func(t *testing.T) {
		resp, err := fetchSocketsFromRemote(mockServer.URL, true, cacheDir, 10, apiHeaderName, apiHeaderValue)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		readBytes, err := io.ReadAll(resp)
		if err != nil {
			t.Fatalf("Failed to read from response: %v\n", err)
		}
		if string(readBytes) != testhelpers.TestSocketList {
			t.Errorf("Expected %s, got %s\n", testhelpers.TestSocketList, string(readBytes))
		}
	})

	t.Run("Fetch with expired cache", func(t *testing.T) {
		resp, err := fetchSocketsFromRemote(mockServer.URL, true, cacheDir, 0, apiHeaderName, apiHeaderValue)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		readBytes, err := io.ReadAll(resp)
		if err != nil {
			t.Fatalf("Failed to read from response: %v\n", err)
		}
		if string(readBytes) != testhelpers.TestSocketList {
			t.Errorf("Expected %s, got %s\n", testhelpers.TestSocketList, string(readBytes))
		}
	})

	t.Run("Fetch without caching", func(t *testing.T) {
		resp, err := fetchSocketsFromRemote(mockServer.URL, false, cacheDir, 0, apiHeaderName, apiHeaderValue)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		readBytes, err := io.ReadAll(resp)
		if err != nil {
			t.Fatalf("Failed to read from response: %v\n", err)
		}
		if string(readBytes) != testhelpers.TestSocketList {
			t.Errorf("Expected %s, got %s\n", testhelpers.TestSocketList, string(readBytes))
		}
	})

	t.Run("Invalid URL without cache flag", func(t *testing.T) {
		badURL := "http://badurl.com"

		_, err := fetchSocketsFromRemote(badURL, false, cacheDir, 0, apiHeaderName, apiHeaderValue)
		if err == nil {
			t.Errorf("Expected error, got none\n")
		}
	})

	t.Run("Invalid URL with cache flag", func(t *testing.T) {
		badURL := "http://badurl.com"

		_, err := fetchSocketsFromRemote(badURL, true, cacheDir, 0, apiHeaderName, apiHeaderValue)
		if err == nil || err == ErrExpiredCache {
			t.Errorf("Expected error, got %v\n", err)
		}
	})
}
