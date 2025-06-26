package socket

import (
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"reflect"
	"testing"

	"go.vxn.dev/dish/pkg/config"
)

func TestNewFetchHandler(t *testing.T) {
	expected := &fetchHandler{
		logger: &mockLogger{},
	}
	actual := NewFetchHandler(&mockLogger{})

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestFetchSocketsFromFile(t *testing.T) {
	filePath := testFile(t, "randomhash.json", []byte(testSockets))
	cfg := &config.Config{
		Source: filePath,
	}

	fetchHandler := NewFetchHandler(&mockLogger{})

	reader, err := fetchHandler.fetchSocketsFromFile(cfg)
	if err != nil {
		t.Fatalf("Failed to fetch sockets from file %v\n", err)
	}
	defer reader.Close()

	fileData, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("Failed to load data from file %v\n", err)
	}

	fileDataString := string(fileData)
	if fileDataString != testSockets {
		t.Errorf("Got %s, expected %s from file\n", fileDataString, testSockets)
	}
}

func TestFetchSocketsFromRemote(t *testing.T) {
	apiHeaderName := "Authorization"
	apiHeaderValue := "Bearer xyzzzzzzz"
	mockServer := newMockServer(t, apiHeaderName, apiHeaderValue, testSockets, http.StatusOK)

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
			filePath := testFile(t, "randomhash.json", []byte(testSockets))
			tt.cfg.ApiCacheDirectory = filepath.Dir(filePath)

			fetchHandler := NewFetchHandler(&mockLogger{})

			resp, err := fetchHandler.fetchSocketsFromRemote(tt.cfg)
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

			if string(readBytes) != testSockets {
				t.Errorf("expected %s, got %s", testSockets, string(readBytes))
			}
		})
	}
}
