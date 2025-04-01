package socket

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// CreateTempFile is a helper function to create a temp file with content
func CreateTempFile(t *testing.T, content string) string {
	t.Helper()

	file, err := os.CreateTemp("", "test_sockets_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer file.Close()

	if _, err := file.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	return file.Name()
}

func TestFetchRemoteStream(t *testing.T) {
	expectedHeaderName := "Authorization"
	expectedHeaderValue := "Bearer [JWT TOKEN]"

	// Simulate valid response
	testServerValid := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaderValue := r.Header.Get(expectedHeaderName)
		if gotHeaderValue != expectedHeaderValue {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer testServerValid.Close()

	// Simulate invalid response
	testServerInvalid := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid request", http.StatusBadRequest)
	}))
	defer testServerValid.Close()

	tests := []struct {
		name        string
		url         string
		headerName  string
		headerValue string
		wantErr     bool
	}{
		{
			name:        "Successful Request",
			url:         testServerValid.URL,
			headerName:  "Authorization",
			headerValue: "Bearer [JWT TOKEN]",
			wantErr:     false,
		},
		{
			name:        "Missing Auth Header",
			url:         testServerValid.URL,
			headerName:  "",
			headerValue: "",
			wantErr:     true,
		},
		{
			name:        "Invalid Request",
			url:         testServerInvalid.URL,
			headerName:  "",
			headerValue: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fetchRemoteStream(tt.url, tt.headerName, tt.headerValue)

			if tt.wantErr && err == nil {
				t.Error("expected an error but got none")
			} else if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestFetchFileStream(t *testing.T) {
	testContent := "testcontent"
	testFile := CreateTempFile(t, testContent)
	defer os.Remove(testFile)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "Valid File Path",
			path:    testFile,
			wantErr: false,
		},
		{
			name:    "Invalid File Path",
			path:    "non_existent.json",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, err := fetchFileStream(tt.path)

			if tt.wantErr {
				if err == nil {
					t.Error("expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			defer reader.Close()
			body, _ := io.ReadAll(reader)
			if string(body) != testContent {
				t.Errorf("expected %s but got %s", testContent, body)
			}
		})
	}
}

func TestGetStreamFromPath(t *testing.T) {
	// This tests only valid sources
	testContent := "testdata"
	testFile := CreateTempFile(t, testContent)
	defer os.Remove(testFile)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testContent))
	}))
	defer testServer.Close()

	tests := []struct {
		name        string
		input       string
		headerName  string
		headerValue string
		expected    string
		wantErr     bool
	}{
		{
			name:     "Path is Filepath",
			input:    testFile,
			expected: testContent,
			wantErr:  false,
		},
		{
			name:     "Path is URL",
			input:    testServer.URL,
			expected: testContent,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, err := getStreamFromPath(tt.input, tt.headerName, tt.headerValue)

			if tt.wantErr {
				if err == nil {
					t.Error("expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			defer reader.Close()
			body, _ := io.ReadAll(reader)
			if string(body) != tt.expected {
				t.Errorf("expected %s but got %s", tt.expected, body)
			}
		})
	}
}
