package socket

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestFetchSocketList(t *testing.T) {
	validJSON := `{
		"sockets": [
			{
				"id": "vxn_dev_https",
				"socket_name": "vxn-dev HTTPS",
				"host_name": "https://vxn.dev",
				"port_tcp": 443,
				"path_http": "/",
				"expected_http_code_array": [200]
			}
		]
	}`
	invalidJSON := `{"sockets": [ { "id": "invalid"`

	validFile := CreateTempFile(t, validJSON)
	defer os.Remove(validFile)

	invalidFile := CreateTempFile(t, invalidJSON)
	defer os.Remove(invalidFile)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(validJSON))
	}))
	defer testServer.Close()

	tests := []struct {
		name      string
		input     string
		wantErr   bool
		expectLen int
	}{
		{
			"Valid File",
			validFile,
			false,
			1,
		},
		{
			"Valid URL",
			testServer.URL,
			false,
			1,
		},
		{
			"Invalid File Path",
			"non_existent.json",
			true,
			0,
		},
		{
			"Malformed JSON",
			invalidFile,
			true,
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := FetchSocketList(tt.input, "", "", false)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(list.Sockets) != tt.expectLen {
				t.Fatalf("expected %d sockets, got %d", tt.expectLen, len(list.Sockets))
			}
		})
	}
}
