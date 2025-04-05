package socket

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
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
	if err := os.MkdirAll(cacheDir, 0o600); err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0o600)
}

// loadCachedSockets checks whether the cache is valid (not expired) and the returns the data stream and ModTime of the cache.
func loadCachedSockets(filePath string, cacheTTL uint) (io.ReadCloser, time.Time, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, time.Time{}, err
	}

	reader, err := os.Open(filePath)
	if err != nil {
		return nil, time.Time{}, err
	}

	cacheTime := info.ModTime()
	if time.Since(cacheTime) > time.Duration(cacheTTL)*time.Minute {
		return reader, cacheTime, ErrExpiredCache
	}

	return reader, cacheTime, nil
}
