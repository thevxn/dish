package socket

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"go.vxn.dev/dish/pkg/testhelpers"
)

func TestHashUrlToFilePath(t *testing.T) {
	tests := []struct {
		url      string
		cacheDir string
		expected string
	}{
		{"https://example.com", "test_cache/", "test_cache/327c3fda87ce286848a574982ddd0b7c7487f816.json"},
		{"http://localhost", "test_cache/", "test_cache/8523ab8065a69338d5006c34310dc8d2c0179ebb.json"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			got := hashUrlToFilePath(tt.url, tt.cacheDir)
			if got != tt.expected {
				t.Errorf("Got %s, want %s\n", got, tt.expected)
			}
		})
	}
}

func TestSaveSocketsToCache(t *testing.T) {
	filePath := testhelpers.TestFile(t, "randomhash.json", nil)
	cacheDir := filepath.Dir(filePath)

	if err := saveSocketsToCache(filePath, cacheDir, []byte(testhelpers.TestSocketList)); err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("Expected file %s to exist, but it doesn't", filePath)
	}

	readBytes, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read saved cache: %v\n", err)
	}

	if string(readBytes) != testhelpers.TestSocketList {
		t.Errorf("Expected file content %s, got %s\n", testhelpers.TestSocketList, string(readBytes))
	}
}

func TestLoadSocketsFromCache(t *testing.T) {
	filePath := testhelpers.TestFile(t, "randomhash.json", []byte(testhelpers.TestSocketList))

	t.Run("Load sockets from cache", func(t *testing.T) {
		cacheTTL := uint(60)
		readerFromCache, err := loadSocketsFromCache(filePath, cacheTTL)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		readBytes, err := io.ReadAll(readerFromCache)
		if err != nil {
			t.Fatalf("Failed to read saved cache: %v\n", err)
		}

		if string(readBytes) != testhelpers.TestSocketList {
			t.Errorf("Expected retrieved data to be %s, got %s", testhelpers.TestSocketList, string(readBytes))
		}
	})

	t.Run("Load sockets from expired cache", func(t *testing.T) {
		cacheTTL := uint(0)
		readerFromCache, err := loadSocketsFromCache(filePath, cacheTTL)
		if !errors.Is(err, ErrExpiredCache) {
			t.Errorf("Expected error %v, but got %v", ErrExpiredCache, err)
		}

		readBytes, err := io.ReadAll(readerFromCache)
		if err != nil {
			t.Fatalf("Failed to read saved cache: %v\n", err)
		}

		if string(readBytes) != testhelpers.TestSocketList {
			t.Errorf("Expected retrieved data to be %s, got %s", testhelpers.TestSocketList, string(readBytes))
		}
	})
}
