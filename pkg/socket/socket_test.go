package socket

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"reflect"
	"testing"

	"go.vxn.dev/dish/pkg/config"
	"go.vxn.dev/dish/pkg/testhelpers"
)

func TestPrintSockets(t *testing.T) {
	list := &SocketList{
		Sockets: []Socket{
			{ID: "1", Name: "socket", Host: "example.com", Port: 80, ExpectedHTTPCodes: []int{200, 404}},
		},
	}

	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)

	PrintSockets(list, logger)

	expected := "Host: example.com, Port: 80, ExpectedHTTPCodes: [200 404]\n"
	if !bytes.Contains(buf.Bytes(), []byte(expected)) {
		t.Errorf("Expected TestPrintSockets() to contain %s, but got %s", expected, buf.String())
	}
}

func TestLoadSocketList(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		expectErr bool
	}{
		{
			"Valid JSON",
			testhelpers.TestSocketList,
			false,
		},
		{
			"Invalid JSON",
			`{ "sockets": [ { "id": "vxn_dev_https"`,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := io.NopCloser(bytes.NewReader([]byte(tt.json)))
			if _, err := LoadSocketList(reader); (err == nil) == tt.expectErr {
				t.Errorf("Expect error: %v, got error: %v\n", tt.expectErr, err)
			}
		})
	}
}

func TestFetchSocketList(t *testing.T) {
	mockServer := testhelpers.NewMockServer(t, "", "", testhelpers.TestSocketList, http.StatusOK)
	validFile := testhelpers.TestFile(t, "randomhash.json", []byte(testhelpers.TestSocketList))
	socketStringReader := io.NopCloser(bytes.NewBufferString(testhelpers.TestSocketList))
	originalList, err := LoadSocketList(socketStringReader)
	if err != nil {
		t.Fatalf("failed to parse sockets string to an object: %v", err)
	}

	newConfig := func(source string) *config.Config {
		return &config.Config{
			Source: source,
			Logger: config.NewLogger(false),
		}
	}

	tests := []struct {
		name        string
		source      string
		expectError bool
	}{
		{
			name:        "Fetch from file",
			source:      validFile,
			expectError: false,
		},
		{
			name:        "Fetch from remote",
			source:      mockServer.URL,
			expectError: false,
		},
		{
			name:        "Fetch from remote with bad URL",
			source:      "http://invalid-host.local",
			expectError: true,
		},
		{
			name:        "Fetch from not existent file",
			source:      "thisdoesntexist.json",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newConfig(tt.source)

			fetchedList, err := FetchSocketList(cfg)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			// Manual comparison of 2 objects won't work because of expected codes type ([]int) in Socket struct
			if !reflect.DeepEqual(fetchedList, originalList) {
				t.Errorf("expected %+v, got %+v", originalList, fetchedList)
			}
		})
	}
}
