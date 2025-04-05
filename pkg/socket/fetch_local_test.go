package socket

import (
	"io"
	"testing"

	"go.vxn.dev/dish/pkg/testhelpers"
)

func TestFetchSocketsFromFile(t *testing.T) {
	filePath := testhelpers.TestFile(t, "randomhash.json", []byte(testhelpers.TestSocketList))

	reader, err := fetchSocketsFromFile(filePath)
	if err != nil {
		t.Fatalf("Failed to fetch sockets from file %v\n", err)
	}
	defer reader.Close()

	fileData, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("Failed to load data from file %v\n", err)
	}

	fileDataString := string(fileData)
	if fileDataString != testhelpers.TestSocketList {
		t.Errorf("Got %s, expected %s from file\n", fileDataString, testhelpers.TestSocketList)
	}
}
