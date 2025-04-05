package socket

import (
	"errors"
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

	t.Run("Fetch With Valid Cache", func(t *testing.T) {
		filePath := testhelpers.TestFile(t, "randomhash.json", []byte(testhelpers.TestSocketList))
		cacheDir := filepath.Dir(filePath)

		resp, err := fetchSocketsFromRemote(mockServer.URL, true, cacheDir, 10, apiHeaderName, apiHeaderValue)
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

	t.Run("Fetch With Expired Cache", func(t *testing.T) {
		filePath := testhelpers.TestFile(t, "randomhash.json", []byte(testhelpers.TestSocketList))
		cacheDir := filepath.Dir(filePath)

		resp, err := fetchSocketsFromRemote(mockServer.URL, true, cacheDir, 0, apiHeaderName, apiHeaderValue)
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

	t.Run("Fetch Without Caching", func(t *testing.T) {
		filePath := testhelpers.TestFile(t, "randomhash.json", []byte(testhelpers.TestSocketList))
		cacheDir := filepath.Dir(filePath)

		resp, err := fetchSocketsFromRemote(mockServer.URL, false, cacheDir, 0, apiHeaderName, apiHeaderValue)
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

	t.Run("Invalid URL Without Cache Flag", func(t *testing.T) {
		filePath := testhelpers.TestFile(t, "randomhash.json", []byte(testhelpers.TestSocketList))
		cacheDir := filepath.Dir(filePath)

		badURL := "http://badurl.com"

		_, err := fetchSocketsFromRemote(badURL, false, cacheDir, 0, apiHeaderName, apiHeaderValue)
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})

	t.Run("Invalid URL With Cache Flag", func(t *testing.T) {
		filePath := testhelpers.TestFile(t, "randomhash.json", []byte(testhelpers.TestSocketList))
		cacheDir := filepath.Dir(filePath)

		badURL := "http://badurl.com"

		_, err := fetchSocketsFromRemote(badURL, true, cacheDir, 0, apiHeaderName, apiHeaderValue)
		if err == nil || errors.Is(err, ErrExpiredCache) {
			t.Errorf("expected error, got %v\n", err)
		}
	})
}
