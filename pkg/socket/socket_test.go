package socket

import (
	"bytes"
	"io"
	"log"
	"testing"

	"go.vxn.dev/dish/pkg/testhelpers"
)

func TestPrintSockets(t *testing.T) {
	list := &SocketList{
		Sockets: []Socket{
			{ID: "1", Name: "socket", Host: "example.com", Port: 80, ExpectedHTTPCodes: []int{200, 404}},
		},
	}

	var buf bytes.Buffer
	log.SetOutput(&buf)

	PrintSockets(list)

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
