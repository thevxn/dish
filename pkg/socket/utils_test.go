package socket

import (
	"testing"
)

func TestIsFilePath(t *testing.T) {
	tests := []struct {
		source     string
		expectFile bool
	}{
		{"path/file.txt", true},
		{"C:/file.txt", true},
		{"./path", true},
		{"file", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run("Test file path", func(t *testing.T) {
			if got := IsFilePath(tt.source); got != tt.expectFile {
				t.Errorf("IsFilePath(%q) = %v, want %v", tt.source, got, tt.expectFile)
			}
		})
	}

	urlTests := []struct {
		source     string
		expectFile bool
	}{
		{"https://example.com", false},
		{"http://localhost:8080", false},
		{"https://www.google.com", false},
	}

	for _, tt := range urlTests {
		t.Run("Test URL", func(t *testing.T) {
			if got := IsFilePath(tt.source); got != tt.expectFile {
				t.Errorf("IsFilePath(%s) = %v, want %v", tt.source, got, tt.expectFile)
			}
		})
	}
}
