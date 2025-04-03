package socket

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// hashUrlToFilePath hashes given URL to create cache file path.
func hashUrlToFilePath(url string, cacheDir string) string {
	hash := sha1.Sum([]byte(url))
	filename := hex.EncodeToString(hash[:]) + ".json"
	return filepath.Join(cacheDir, filename)
}

// saveSocketsToCache caches sockets to specified file in cache directory.
func saveSocketsToCache(filePath string, cacheDir string, reader io.ReadCloser) error {
	defer reader.Close()

	// Make sure that cache directory exists
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}

// loadSocketsFromCache checks if cache is valid (TTL left) and returns data stream.
func loadSocketsFromCache(filePath string, cacheTTL uint) (io.ReadCloser, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	} else if time.Since(info.ModTime()) > time.Duration(cacheTTL)*time.Hour {
		return nil, fmt.Errorf("cache for this source is outdated (%s)", filePath)
	}

	return os.Open(filePath)
}
