package socket

import (
	"io"
	"os"
)

// fetchSocketsFromFile opens a file and returns [io.ReadCloser] for reading from the stream.
func fetchSocketsFromFile(input string) (io.ReadCloser, error) {
	return os.Open(input)
}
