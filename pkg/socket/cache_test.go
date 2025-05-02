package socket

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.vxn.dev/dish/pkg/testhelpers"
)

func TestHashUrlToFilePath(t *testing.T) {
	tests := []struct {
		url      string
		cacheDir string
		expected string
	}{
		{
			"https://example.com",
			"test_cache",
			filepath.Join("test_cache", "327c3fda87ce286848a574982ddd0b7c7487f816.json"),
		},
		{
			"http://localhost",
			"test_cache",
			filepath.Join("test_cache", "8523ab8065a69338d5006c34310dc8d2c0179ebb.json"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			got := hashUrlToFilePath(tt.url, tt.cacheDir)
			if got != tt.expected {
				t.Errorf("got %s, want %s", got, tt.expected)
			}
		})
	}
}

func TestSaveSocketsToCache(t *testing.T) {
	filePath := testhelpers.TestFile(t, "randomhash.json", nil)
	cacheDir := filepath.Dir(filePath)

	if err := saveSocketsToCache(filePath, cacheDir, []byte(testhelpers.TestSocketList)); err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("expected file %s to exist, but it does not", filePath)
	}

	readBytes, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read saved cache: %v", err)
	}

	if string(readBytes) != testhelpers.TestSocketList {
		t.Errorf("expected file content %s, got %s", testhelpers.TestSocketList, string(readBytes))
	}
}

func TestLoadSocketsFromCache(t *testing.T) {
	t.Run("Load Sockets From Cache", func(t *testing.T) {
		filePath := testhelpers.TestFile(t, "randomhash.json", []byte(testhelpers.TestSocketList))
		cacheTTL := uint(60)

		readerFromCache, _, err := loadCachedSockets(filePath, cacheTTL)
		if err != nil {
			t.Fatalf("expected no error, but got %v", err)
		}
		defer readerFromCache.Close()

		readBytes, err := io.ReadAll(readerFromCache)
		if err != nil {
			t.Fatalf("failed to read saved cache: %v", err)
		}

		if string(readBytes) != testhelpers.TestSocketList {
			t.Errorf("expected retrieved data to be %s, got %s", testhelpers.TestSocketList, string(readBytes))
		}
	})

	t.Run("Load Sockets From Expired Cache", func(t *testing.T) {
		filePath := testhelpers.TestFile(t, "randomhash.json", []byte(testhelpers.TestSocketList))
		cacheTTL := uint(0)
		time.Sleep(200 * time.Millisecond)

		readerFromCache, _, err := loadCachedSockets(filePath, cacheTTL)
		if !errors.Is(err, ErrExpiredCache) {
			t.Errorf("expected error %v, but got %v", ErrExpiredCache, err)
		}
		defer readerFromCache.Close()

		readBytes, err := io.ReadAll(readerFromCache)
		if err != nil {
			t.Fatalf("failed to read saved cache: %v", err)
		}

		if string(readBytes) != testhelpers.TestSocketList {
			t.Errorf("expected retrieved data to be %s, got %s", testhelpers.TestSocketList, string(readBytes))
		}
	})
}
