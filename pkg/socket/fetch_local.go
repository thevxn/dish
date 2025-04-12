package socket

import (
	"io"
	"log"
	"os"

	"go.vxn.dev/dish/pkg/config"
)

// fetchSocketsFromFile opens a file and returns [io.ReadCloser] for reading from the stream.
func fetchSocketsFromFile(config *config.Config) (io.ReadCloser, error) {
	file, err := os.Open(config.Source)
	if err != nil {
		return nil, err
	}

	// TODO: Replace with logger
	if config.Verbose {
		log.Printf("fetching sockets from file (%s)", config.Source)
	}

	return file, nil
}
