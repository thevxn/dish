package socket

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

var ErrExpiredCache error = errors.New("cache file for this source is outdated")

// hashUrlToFilePath hashes given URL to create cache file path.
func hashUrlToFilePath(url string, cacheDir string) string {
	hash := sha1.Sum([]byte(url))
	filename := hex.EncodeToString(hash[:]) + ".json"
	return filepath.Join(cacheDir, filename)
}

// saveSocketsToCache caches socket data to specified file in cache directory.
func saveSocketsToCache(filePath string, cacheDir string, data []byte) error {
	// Make sure that cache directory exists
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	n, err := file.Write(data)
	if err != nil {
		return err
	}

	if n < len(data) {
		return fmt.Errorf("incomplete write: wrote %d/%d bytes", n, len(data))
	}

	return nil
}

// loadSocketsFromCache checks if cache is not expired and returns data stream.
func loadSocketsFromCache(filePath string, cacheTTL uint) (io.ReadCloser, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	reader, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	if time.Since(info.ModTime()) > time.Duration(cacheTTL)*time.Minute {
		return reader, ErrExpiredCache
	}

	return reader, nil
}
