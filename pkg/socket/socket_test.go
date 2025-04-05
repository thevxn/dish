package socket

import (
	"bytes"
	"io"
	"log"
	"net/http"
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

func TestFetchSocketList(t *testing.T) {
	t.Run("Fetch from file", func(t *testing.T) {
		path := testhelpers.TestFile(t, "randomhash.json", []byte(testhelpers.TestSocketList))

		list, err := FetchSocketList(path, false, "", 0, "", "")
		if err != nil {
			t.Fatal(err)
		}

		if len(list.Sockets) != 1 {
			t.Errorf("Expected list length to be 1, got %d elements\n", len(list.Sockets))
		}

		expectedID := "vxn_dev_https"
		if expectedID != list.Sockets[0].ID {
			t.Errorf("Expected ID=%s, got ID=%s\n", expectedID, list.Sockets[0].ID)
		}
	})

	t.Run("Fetch from remote", func(t *testing.T) {
		server := testhelpers.NewMockServer(t, "", "", testhelpers.TestSocketList, http.StatusOK)

		list, err := FetchSocketList(server.URL, false, "", 0, "", "")
		if err != nil {
			t.Fatal(err)
		}

		if len(list.Sockets) != 1 {
			t.Errorf("Expected list length to be 1, got %d elements\n", len(list.Sockets))
		}

		expectedID := "vxn_dev_https"
		if expectedID != list.Sockets[0].ID {
			t.Errorf("Expected ID=%s, got ID=%s\n", expectedID, list.Sockets[0].ID)
		}
	})

	t.Run("Fetch from remote with bad URL", func(t *testing.T) {
		_, err := FetchSocketList("http://invalid-host.local", false, "", 0, "", "")
		if err == nil {
			t.Errorf("Expected an error got nil\n")
		}
	})

	t.Run("Fetch from not existent file", func(t *testing.T) {
		_, err := FetchSocketList("thisdoesnotexist.json", false, "", 0, "", "")
		if err == nil {
			t.Errorf("Expected an error got nil\n")
		}
	})
}
