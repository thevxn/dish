package testhelpers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// This socket list is used across tests.
const TestSocketList string = `{ "sockets": [ { "id": "vxn_dev_https", "socket_name": "vxn-dev HTTPS", "host_name": "https://vxn.dev", "port_tcp": 443, "path_http": "/", "expected_http_code_array": [200] } ] }`

// TestFile creates a temporary file with provided filename and data.
// Directory is removed when test finishes.
func TestFile(t *testing.T, filename string, data []byte) string {
	t.Helper()
	dir, err := os.MkdirTemp(t.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}

	filepath := filepath.Join(dir, filename)

	err = os.WriteFile(filepath, data, 0o600)
	if err != nil {
		t.Fatal(err)
	}

	return filepath
}

// NewMockServer creates an httptest.Server that simulates an expected API endpoint.
// It validates a specific request header (if provided) and returns a customizable response.
func NewMockServer(t *testing.T, expectedHeaderName, expectedHeaderValue, responseBody string, statusCode int) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if expectedHeaderName != "" && expectedHeaderValue != "" {
			if r.Header.Get(expectedHeaderName) != expectedHeaderValue {
				http.Error(w, `{"error":"Invalid or missing header"}`, http.StatusForbidden)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write([]byte(responseBody))
	}))

	// Automatically close the server after the test finishes
	t.Cleanup(func() {
		server.Close()
	})

	return server
}
