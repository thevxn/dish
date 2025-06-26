package main

import (
	"os"
	"path/filepath"
	"testing"
)

// This socket list is used across tests. It contains a socket which should always pass checks (unless the site is actually down).
const testSocketsValid string = `{ "sockets": [ { "id": "vxn_dev_https", "socket_name": "vxn-dev HTTPS", "host_name": "https://vxn.dev", "port_tcp": 443, "path_http": "/", "expected_http_code_array": [200] } ] }`

// This socket list is used across tests. It contains a socket which should never pass checks.
const testSocketsSomeInvalid string = `{ "sockets": [ { "id": "vxn_dev_https", "socket_name": "vxn-dev HTTPS", "host_name": "https://vxn.dev", "port_tcp": 443, "path_http": "/", "expected_http_code_array": [500] } ] }`

// testFile creates a temporary file inside of a temporary directory with the provided filename and data.
// The temporary directory including the file is removed when the test using it finishes.
func testFile(t *testing.T, filename string, data []byte) string {
	t.Helper()
	dir := t.TempDir()

	filepath := filepath.Join(dir, filename)

	err := os.WriteFile(filepath, data, 0o600)
	if err != nil {
		t.Fatal(err)
	}

	return filepath
}
